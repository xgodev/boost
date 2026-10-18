package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/builder/examples/simple/gen/builder/api"
	grpcserver "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
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
	)
	app := fx.New(append(api.Modules(),
		fx.NopLogger,
		fx.Populate(&rest, &grpcSrv),
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
}
