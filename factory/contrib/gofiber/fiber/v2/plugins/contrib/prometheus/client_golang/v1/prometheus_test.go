package prometheus

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	f "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// newServer builds a server from the config defaults with plugins and a GET
// /x route.
func newServer(t *testing.T, plugins ...fiber.Plugin) *fiber.Server {
	t.Helper()
	opts, err := fiber.NewOptions()
	require.NoError(t, err)
	srv, err := fiber.NewServerWithOptions(context.Background(), opts, plugins...)
	require.NoError(t, err)
	srv.Get("/x", func(c *f.Ctx) error { return c.SendString("x") })
	return srv
}

// do answers req through the app and returns the response and its body.
func do(t *testing.T, srv *fiber.Server, req *http.Request) (*http.Response, string) {
	t.Helper()
	res, err := srv.App().Test(req, -1)
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return res, string(body)
}

func TestDisabledAddsNothing(t *testing.T) {
	srv := newServer(t, NewPrometheusWithOptions(&Options{}).Register)
	assert.Equal(t, newServer(t).App().HandlersCount(), srv.App().HandlersCount())
}

func TestRegisterCountsRequestsByRouteAndServesMetrics(t *testing.T) {
	srv := newServer(t, Register)
	srv.Get("/orders/:id", func(c *f.Ctx) error { return c.SendString("ok") })

	res, _ := do(t, srv, httptest.NewRequest(http.MethodGet, "/orders/7", nil))
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, body := do(t, srv, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Contains(t, body, `fiber_requests_total{code="200",method="GET",url="/orders/:id"}`)
	assert.Contains(t, body, `fiber_request_duration_seconds_count{code="200",method="GET",url="/orders/:id"}`)
}

func TestRegisterOnASecondServerSharesTheCollectors(t *testing.T) {
	newServer(t, Register)

	opts, err := fiber.NewOptions()
	require.NoError(t, err)
	_, err = fiber.NewServerWithOptions(context.Background(), opts, Register)

	assert.NoError(t, err, "a second server must not register the fiber collectors again")
}

func TestRegisterWithoutTheRouteKeepsTheMiddleware(t *testing.T) {
	o, err := NewOptions()
	require.NoError(t, err)
	o.Route.Enabled = false

	srv := newServer(t, NewPrometheusWithOptions(o).Register)

	res, _ := do(t, srv, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Greater(t, srv.App().HandlersCount(), newServer(t).App().HandlersCount())
}
