package error_handler

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
	"github.com/xgodev/boost/model/errors"
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
	srv := newServer(t, NewErrorHandlerWithOptions(&Options{}).Register)
	assert.Equal(t, newServer(t).App().HandlersCount(), srv.App().HandlersCount())
}

func TestRegisterAnswersErrorsAsJSON(t *testing.T) {
	srv := newServer(t, Register)
	srv.Get("/err", func(*f.Ctx) error { return errors.NotFoundf("order 1") })

	res, body := do(t, srv, httptest.NewRequest(http.MethodGet, "/err", nil))

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Contains(t, body, `"message":"order 1 not found"`)
}

func TestRegisterAnswersPlainTextOutsideREST(t *testing.T) {
	opts, err := fiber.NewOptions()
	require.NoError(t, err)
	opts.Type = "HTML"
	srv, err := fiber.NewServerWithOptions(context.Background(), opts, Register)
	require.NoError(t, err)
	srv.Get("/err", func(*f.Ctx) error { return errors.NotFoundf("order 1") })

	res, body := do(t, srv, httptest.NewRequest(http.MethodGet, "/err", nil))

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "order 1 not found", body)
}
