package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	natsserver "github.com/nats-io/nats-server/v2/test"
	natsgo "github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	stockv1 "github.com/xgodev/boost/examples/stock/api/proto/stock/v1"
	"github.com/xgodev/boost/examples/stock/gen/builder/api"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/reserver"
	grpcserver "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	mcpserver "github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// TestMain boots boost: the servers read their config.
func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func serve(srv *boostecho.Server, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, req)
	return rec
}

// TestAPI builds app api from the generated Modules() once (the modules
// register themselves once per process) and calls each door without
// listening on a port. Each door works on its own items of the one memory
// stock: rest on tshirt and cap, grpc on pen, mcp on mug.
func TestAPI(t *testing.T) {
	var (
		rest    *boostecho.Server
		grpcSrv *grpcserver.Server
		mcpSrv  *mcpserver.Server
		reserve reserver.Service
	)
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&rest, &grpcSrv, &mcpSrv, &reserve))...)
	require.NoError(t, app.Err())

	t.Run("rest creates, finds and reserves", func(t *testing.T) {
		got := serve(rest, http.MethodGet, "/stocks/tshirt", "")
		assert.Equal(t, http.StatusOK, got.Code)
		assert.JSONEq(t, `{"sku":"tshirt","units":3}`, got.Body.String(), "the memory stock starts with its seed")

		created := serve(rest, http.MethodPost, "/stocks", `{"sku":"cap","units":1}`)
		assert.Equal(t, http.StatusCreated, created.Code)
		assert.JSONEq(t, `{"sku":"cap","units":1}`, created.Body.String())

		reserved := serve(rest, http.MethodPost, "/stocks/cap/reserve", "")
		assert.Equal(t, http.StatusOK, reserved.Code)
		assert.JSONEq(t, `{"sku":"cap","units":0}`, reserved.Body.String(), "the reservation answers what is left")
	})

	t.Run("rest answers each error kind with its status", func(t *testing.T) {
		for _, c := range []struct {
			name, method, target, body string
			status                     int
		}{
			{"unknown sku", http.MethodGet, "/stocks/nope", "", http.StatusNotFound},
			{"blank sku", http.MethodPost, "/stocks", `{"sku":"","units":1}`, http.StatusBadRequest},
			{"negative units", http.MethodPost, "/stocks", `{"sku":"hat","units":-1}`, http.StatusBadRequest},
			{"body that is not json", http.MethodPost, "/stocks", `{`, http.StatusBadRequest},
			{"sku already in stock", http.MethodPost, "/stocks", `{"sku":"tshirt","units":1}`, http.StatusConflict},
			{"sold out", http.MethodPost, "/stocks/cap/reserve", "", http.StatusConflict},
		} {
			t.Run(c.name, func(t *testing.T) {
				assert.Equal(t, c.status, serve(rest, c.method, c.target, c.body).Code)
			})
		}
	})

	t.Run("rest serves health and metrics", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, serve(rest, http.MethodGet, "/health", "").Code)
		assert.Equal(t, http.StatusOK, serve(rest, http.MethodGet, "/metrics", "").Code)
	})

	t.Run("the services run in their middlewares", func(t *testing.T) {
		assert.Contains(t, reflect.TypeOf(reserve).Elem().PkgPath(), "gen/builder/middleware/",
			"every consumer receives the decorator builder generates for @Middleware")
		assert.GreaterOrEqual(t, opsTotal(t, "reserver.service.reserve"), 1.0, "the prometheus middleware counted the reservations")
		assert.GreaterOrEqual(t, opsTotal(t, "finder.service.find"), 1.0)
		assert.GreaterOrEqual(t, opsTotal(t, "creator.service.create"), 1.0)
	})

	t.Run("grpc finds and reserves", func(t *testing.T) {
		require.Equal(t, http.StatusCreated, serve(rest, http.MethodPost, "/stocks", `{"sku":"pen","units":1}`).Code)

		lis := bufconn.Listen(1 << 20)
		go func() { _ = grpcSrv.Server().Serve(lis) }()
		t.Cleanup(grpcSrv.Server().Stop)

		conn, err := grpc.NewClient("passthrough:///bufnet",
			grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		t.Cleanup(func() { _ = conn.Close() })

		ctx := context.Background()
		found, err := stockv1.NewStockFinderClient(conn).Find(ctx, &stockv1.FindRequest{Sku: "pen"})
		require.NoError(t, err)
		assert.Equal(t, int64(1), found.GetUnits())

		left, err := stockv1.NewStockReserverClient(conn).Reserve(ctx, &stockv1.ReserveRequest{Sku: "pen"})
		require.NoError(t, err)
		assert.Equal(t, int64(0), left.GetUnits())

		_, err = stockv1.NewStockReserverClient(conn).Reserve(ctx, &stockv1.ReserveRequest{Sku: "pen"})
		assert.Equal(t, codes.AlreadyExists, status.Code(err), "boost's grpc server answers a conflict with AlreadyExists")

		_, err = stockv1.NewStockFinderClient(conn).Find(ctx, &stockv1.FindRequest{Sku: "nope"})
		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("mcp serves the find_stock tool", func(t *testing.T) {
		ctx := context.Background()
		serverTransport, clientTransport := mcp.NewInMemoryTransports()
		serverSession, err := mcpSrv.Instance().Connect(ctx, serverTransport, nil)
		require.NoError(t, err)
		t.Cleanup(func() { _ = serverSession.Close() })

		session, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil).Connect(ctx, clientTransport, nil)
		require.NoError(t, err)
		t.Cleanup(func() { _ = session.Close() })

		tools, err := session.ListTools(ctx, nil)
		require.NoError(t, err)
		require.Len(t, tools.Tools, 1)
		assert.Equal(t, "find_stock", tools.Tools[0].Name)
		assert.Equal(t, "Tells how many units of an item are in stock.", tools.Tools[0].Description, "the description is the handler's doc comment")

		res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "find_stock", Arguments: map[string]any{"sku": "mug"}})
		require.NoError(t, err)
		assert.False(t, res.IsError)
		assert.Equal(t, map[string]any{"sku": "mug", "units": 1.0}, res.StructuredContent)

		res, err = session.CallTool(ctx, &mcp.CallToolParams{Name: "find_stock", Arguments: map[string]any{"sku": "nope"}})
		require.NoError(t, err)
		assert.True(t, res.IsError, "an unknown sku is a tool error")
	})
}

const childEnv = "STOCK_EXAMPLE_CHILD"

// child runs test in a new process with env set: config is read once per
// process.
func child(t *testing.T, test string, env ...string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^"+test+"$", "-test.v")
	cmd.Env = append(append(os.Environ(), childEnv+"=1"), env...)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", out)
	assert.Contains(t, string(out), "--- PASS: "+test)
}

// seedEnv sets the items the memory stock starts with
// (internal.pkg.datastore.stock.memory.seed, from its @Config Options).
const seedEnv = "INTERNAL_PKG_DATASTORE_STOCK_MEMORY_SEED"

// TestAPIConfig boots app api with the seed of the memory stock set from the
// environment.
func TestAPIConfig(t *testing.T) {
	child(t, "TestAPIConfigChild", seedEnv+"=hat=5")
}

// TestAPIConfigChild runs in the process TestAPIConfig starts.
func TestAPIConfigChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestAPIConfig")
	}

	var rest *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&rest))...)
	require.NoError(t, app.Err())

	assert.JSONEq(t, `{"sku":"hat","units":5}`, serve(rest, http.MethodGet, "/stocks/hat", "").Body.String())
	assert.Equal(t, http.StatusNotFound, serve(rest, http.MethodGet, "/stocks/tshirt", "").Code, "the env replaces the default seed")
}

const (
	// driverEnv chooses the driver of the publisher that announces the
	// reservations (boost.wrapper.publisher.driver.name, noop by default).
	driverEnv  = "BOOST_WRAPPER_PUBLISHER_DRIVER_NAME"
	natsURLEnv = "BOOST_FACTORY_NATS_URL"
)

// TestAPIPublisher boots app api with the nats driver on a NATS server this
// test runs: a reservation is announced as a stock.reserved CloudEvent.
func TestAPIPublisher(t *testing.T) {
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	srv := natsserver.RunServer(&opts)
	t.Cleanup(srv.Shutdown)

	child(t, "TestAPIPublisherChild", driverEnv+"=nats", natsURLEnv+"="+srv.ClientURL())
}

// TestAPIPublisherChild runs in the process TestAPIPublisher starts.
func TestAPIPublisherChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestAPIPublisher")
	}

	conn, err := natsgo.Connect(os.Getenv(natsURLEnv))
	require.NoError(t, err)
	t.Cleanup(conn.Close)
	sub, err := conn.SubscribeSync("stock.reserved")
	require.NoError(t, err)

	var rest *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&rest))...)
	require.NoError(t, app.Err())
	require.Equal(t, http.StatusOK, serve(rest, http.MethodPost, "/stocks/tshirt/reserve", "").Code)

	msg, err := sub.NextMsg(5 * time.Second)
	require.NoError(t, err, "the reservation was not announced on stock.reserved")
	var event struct {
		Type string `json:"type"`
		Data struct {
			SKU   string `json:"sku"`
			Units int64  `json:"units"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(msg.Data, &event), "%s", msg.Data)
	assert.Equal(t, "stock.reserved", event.Type)
	assert.Equal(t, "tshirt", event.Data.SKU)
	assert.Equal(t, int64(2), event.Data.Units)
}

const (
	// repositoryEnv chooses the implementation of the stock repository
	// (internal.pkg.domain.repository.stock.implementation, memory by default).
	repositoryEnv = "INTERNAL_PKG_DOMAIN_REPOSITORY_STOCK_IMPLEMENTATION"
	mongoURIEnv   = "BOOST_FACTORY_MONGO_URI"
)

// TestAPIMongo boots app api on the mongodb stock. It needs a MongoDB:
// BOOST_FACTORY_MONGO_URI=mongodb://localhost:27017/stock go test ./...
func TestAPIMongo(t *testing.T) {
	if os.Getenv(mongoURIEnv) == "" {
		t.Skip(mongoURIEnv + " is not set")
	}
	child(t, "TestAPIMongoChild", repositoryEnv+"=mongodb",
		"INTERNAL_PKG_DATASTORE_STOCK_MONGODB_COLLECTION=stock_"+time.Now().Format("20060102150405.000000000"))
}

// TestAPIMongoChild runs in the process TestAPIMongo starts, on a collection
// of its own.
func TestAPIMongoChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestAPIMongo")
	}

	var rest *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&rest))...)
	require.NoError(t, app.Err())

	assert.Equal(t, http.StatusNotFound, serve(rest, http.MethodGet, "/stocks/cap", "").Code)
	assert.Equal(t, http.StatusCreated, serve(rest, http.MethodPost, "/stocks", `{"sku":"cap","units":1}`).Code)
	assert.Equal(t, http.StatusConflict, serve(rest, http.MethodPost, "/stocks", `{"sku":"cap","units":1}`).Code)
	assert.JSONEq(t, `{"sku":"cap","units":0}`, serve(rest, http.MethodPost, "/stocks/cap/reserve", "").Body.String())
	assert.Equal(t, http.StatusConflict, serve(rest, http.MethodPost, "/stocks/cap/reserve", "").Code)
	assert.Equal(t, http.StatusNotFound, serve(rest, http.MethodPost, "/stocks/nope/reserve", "").Code)
}

// opsTotal returns the ops_total the prometheus middleware counted for the
// wrapper name.
func opsTotal(t *testing.T, name string) float64 {
	t.Helper()
	families, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)
	for _, family := range families {
		if family.GetName() != "ops_total" {
			continue
		}
		for _, metric := range family.GetMetric() {
			for _, label := range metric.GetLabel() {
				if label.GetName() == "name" && label.GetValue() == name {
					return metric.GetCounter().GetValue()
				}
			}
		}
	}
	return 0
}
