package echoprometheus

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	e "github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func newServer(t *testing.T) *echo.Server {
	t.Helper()
	srv, err := echo.NewServerWithOptions(context.Background(), &echo.Options{Type: "REST"})
	require.NoError(t, err)
	return srv
}

func do(srv *echo.Server, method, target string, body io.Reader) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(method, target, body))
	return rec
}

func ping(srv *echo.Server) {
	srv.GET("/ping", func(c *e.Context) error { return c.String(http.StatusOK, "pong") })
}

func TestRegisterCountsRequestsAndServesMetrics(t *testing.T) {
	srv := newServer(t)
	require.NoError(t, Register(context.Background(), srv))
	ping(srv)

	require.Equal(t, http.StatusOK, do(srv, http.MethodGet, "/ping", nil).Code)

	rec := do(srv, http.MethodGet, "/metrics", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `echo_requests_total{code="200",host="example.com",method="GET",url="/ping"}`)
}

func TestRegisterOnASecondServerSharesTheCollectors(t *testing.T) {
	require.NoError(t, Register(context.Background(), newServer(t)))

	assert.NoError(t, Register(context.Background(), newServer(t)),
		"a second server must not register the echo collectors again")
}

func TestRegisterWithoutTheRouteKeepsTheMiddleware(t *testing.T) {
	srv := newServer(t)
	o, err := NewOptions()
	require.NoError(t, err)
	o.Route.Enabled = false

	require.NoError(t, NewPrometheusWithOptions(o).Register(context.Background(), srv))

	assert.Equal(t, http.StatusNotFound, do(srv, http.MethodGet, "/metrics", nil).Code)
}

func TestRegisterDisabledDoesNothing(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, NewPrometheusWithOptions(&Options{}).Register(context.Background(), srv))

	assert.Equal(t, http.StatusNotFound, do(srv, http.MethodGet, "/metrics", nil).Code)
}
