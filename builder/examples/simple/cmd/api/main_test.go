package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/builder/examples/simple/gen/builder/api"
	"github.com/xgodev/boost/builder/examples/simple/stock"
	grpcserver "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	mcpserver "github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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
// listening on a port.
func TestAPI(t *testing.T) {
	var (
		rest, admin *boostecho.Server
		grpcSrv     *grpcserver.Server
		mcpSrv      *mcpserver.Server
		stockSvc    stock.Service
	)
	app := fx.New(append(api.Modules(),
		fx.NopLogger,
		fx.Populate(&rest, &grpcSrv, &mcpSrv, &stockSvc),
		fx.Invoke(fx.Annotate(func(srv *boostecho.Server) { admin = srv }, fx.ParamTags(`name:"admin"`))),
	)...)
	require.NoError(t, app.Err())
	require.NotSame(t, rest, admin, "the admin server is its own echo instance")

	t.Run("rest serves the order routes", func(t *testing.T) {
		created := serve(rest, http.MethodPost, "/orders", `{"item":"book"}`)
		assert.Equal(t, http.StatusCreated, created.Code)
		assert.JSONEq(t, `{"id":"1","item":"book"}`, created.Body.String())
		assert.Empty(t, created.Header().Get("X-Served-By"), "the admin plugins are admin's only")

		got := serve(rest, http.MethodGet, "/orders/1", "")
		assert.Equal(t, http.StatusOK, got.Code)
		assert.JSONEq(t, `{"id":"1","item":"book"}`, got.Body.String())

		assert.Equal(t, http.StatusNotFound, serve(rest, http.MethodGet, "/orders/9", "").Code)
		assert.Equal(t, http.StatusNotFound, serve(rest, http.MethodGet, "/live", "").Code, "the probe is on admin only")
	})

	t.Run("the stock service runs in its middlewares", func(t *testing.T) {
		assert.Contains(t, reflect.TypeOf(stockSvc).Elem().PkgPath(), "gen/builder/middleware/",
			"every consumer receives the decorator builder generates for @Middleware")

		conflict := serve(rest, http.MethodPost, "/orders", `{"item":"sold-out"}`)
		assert.Equal(t, http.StatusConflict, conflict.Code, "the handler calls the real service through the decorator")
		assert.GreaterOrEqual(t, opsTotal(t, "stock.service.available"), 1.0, "the prometheus middleware counted the call")
	})

	t.Run("admin serves the probe only", func(t *testing.T) {
		live := serve(admin, http.MethodGet, "/live", "")
		assert.Equal(t, http.StatusOK, live.Code)
		assert.Equal(t, "ok", live.Body.String())
		assert.Equal(t, "admin", live.Header().Get("X-Served-By"), "admin is built with the plugins provided under its name")

		assert.Equal(t, http.StatusOK, serve(admin, http.MethodHead, "/live", "").Code)
		assert.Equal(t, http.StatusNotFound, serve(admin, http.MethodGet, "/orders/1", "").Code)
	})

	t.Run("grpc serves the health service", func(t *testing.T) {
		lis := bufconn.Listen(1 << 20)
		go func() { _ = grpcSrv.Server().Serve(lis) }()
		t.Cleanup(grpcSrv.Server().Stop)

		conn, err := grpc.NewClient("passthrough:///bufnet",
			grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		t.Cleanup(func() { _ = conn.Close() })

		resp, err := healthpb.NewHealthClient(conn).Check(context.Background(), &healthpb.HealthCheckRequest{Service: "orders"})
		require.NoError(t, err)
		assert.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
	})

	t.Run("mcp serves the stock tool", func(t *testing.T) {
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
		assert.Equal(t, "check_stock", tools.Tools[0].Name)
		assert.Equal(t, "Tells whether an item can be ordered.", tools.Tools[0].Description, "the description is the handler's doc comment")

		for item, available := range map[string]bool{"book": true, "sold-out": false} {
			res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "check_stock", Arguments: map[string]any{"item": item}})
			require.NoError(t, err)
			assert.False(t, res.IsError)
			assert.Equal(t, map[string]any{"available": available}, res.StructuredContent, "the tool calls the stock service")
		}
	})
}

const (
	// implementationEnv chooses the implementation of stock.Service
	// (builder.examples.simple.stock.implementation).
	implementationEnv = "BUILDER_EXAMPLES_SIMPLE_STOCK_IMPLEMENTATION"
	childEnv          = "BUILDER_EXAMPLE_CHILD"
	// warehouseHostEnv points the warehouse client builder generates for
	// stock/remote at a server (boost.factory.resty.clients.warehouse.host).
	warehouseHostEnv = "BOOST_FACTORY_RESTY_CLIENTS_WAREHOUSE_HOST"
)

// TestAPIImplementation boots app api in a new process per value of the
// implementation key: config is read once per process. The remote
// implementation asks a warehouse this test serves.
func TestAPIImplementation(t *testing.T) {
	warehouse := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || !strings.HasPrefix(r.URL.Path, "/items/") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		_, _ = w.Write([]byte(`{"available":` + strconv.FormatBool(r.URL.Path != "/items/sold-out") + `}`))
	}))
	t.Cleanup(warehouse.Close)

	for _, value := range []string{"unlimited", "remote", "nope"} {
		t.Run(value, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestAPIImplementationChild$", "-test.v")
			cmd.Env = append(os.Environ(), implementationEnv+"="+value, warehouseHostEnv+"="+warehouse.URL, childEnv+"=1")
			out, err := cmd.CombinedOutput()
			require.NoError(t, err, "%s", out)
			assert.Contains(t, string(out), "--- PASS: TestAPIImplementationChild")
		})
	}
}

// TestAPIImplementationChild runs in the process TestAPIImplementation
// starts, with the implementation key set.
func TestAPIImplementationChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestAPIImplementation")
	}

	var (
		rest     *boostecho.Server
		stockSvc stock.Service
	)
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&rest, &stockSvc))...)

	switch value := os.Getenv(implementationEnv); value {
	case "unlimited":
		require.NoError(t, app.Err())
		assert.True(t, strings.HasSuffix(reflect.TypeOf(stockSvc).Elem().PkgPath(), "/stock/unlimited"),
			"config wires unlimited, undecorated: the @Middleware of memory is off with memory")
		assert.Equal(t, http.StatusCreated, serve(rest, http.MethodPost, "/orders", `{"item":"sold-out"}`).Code,
			"unlimited has every item")
	case "remote":
		require.NoError(t, app.Err())
		assert.True(t, strings.HasSuffix(reflect.TypeOf(stockSvc).Elem().PkgPath(), "/stock/remote"),
			"config wires remote, which asks the warehouse through the client builder generates")
		assert.Equal(t, http.StatusCreated, serve(rest, http.MethodPost, "/orders", `{"item":"book"}`).Code,
			"the warehouse has the book")
		assert.Equal(t, http.StatusConflict, serve(rest, http.MethodPost, "/orders", `{"item":"sold-out"}`).Code,
			"the warehouse is out of sold-out")
	default:
		require.Error(t, app.Err(), "an unknown implementation fails the boot")
		assert.Contains(t, app.Err().Error(), "builder.examples.simple.stock.implementation")
		assert.Contains(t, app.Err().Error(), `"`+value+`"`)
		assert.Contains(t, app.Err().Error(), "memory, remote, unlimited")
	}
}

// soldOutEnv sets the items the memory stock is out of
// (builder.examples.simple.stock.memory.soldOut, from its @Config Options).
const soldOutEnv = "BUILDER_EXAMPLES_SIMPLE_STOCK_MEMORY_SOLD__OUT"

// TestAPIConfig boots app api in a new process with the sold-out list of
// the memory stock set from the environment: config is read once per
// process.
func TestAPIConfig(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestAPIConfigChild$", "-test.v")
	cmd.Env = append(os.Environ(), soldOutEnv+"=book,pen", childEnv+"=1")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", out)
	assert.Contains(t, string(out), "--- PASS: TestAPIConfigChild")
}

// TestAPIConfigChild runs in the process TestAPIConfig starts.
func TestAPIConfigChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestAPIConfig")
	}

	var rest *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&rest))...)
	require.NoError(t, app.Err())

	for _, item := range []string{"book", "pen"} {
		assert.Equal(t, http.StatusConflict, serve(rest, http.MethodPost, "/orders", `{"item":"`+item+`"}`).Code,
			"%s is in the sold-out list set from %s", item, soldOutEnv)
	}
	assert.Equal(t, http.StatusCreated, serve(rest, http.MethodPost, "/orders", `{"item":"sold-out"}`).Code,
		"the env replaces the default list")
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
