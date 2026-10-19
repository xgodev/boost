package log

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
	blog "github.com/xgodev/boost/wrapper/log"
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
	srv := newServer(t, NewLogWithOptions(&Options{}).Register)
	assert.Equal(t, newServer(t).App().HandlersCount(), srv.App().HandlersCount())
}

func TestRegisterPutsTheLoggerInTheUserContext(t *testing.T) {
	srv := newServer(t, Register)
	srv.Get("/ctx", func(c *f.Ctx) error {
		assert.NotNil(t, blog.FromContext(c.UserContext()))
		return c.SendString("ok")
	})

	res, body := do(t, srv, httptest.NewRequest(http.MethodGet, "/ctx", nil))

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "ok", body)
}

func TestRegisterAnswersHandlerErrorsThroughTheErrorHandler(t *testing.T) {
	srv := newServer(t, Register)
	srv.Get("/err", func(*f.Ctx) error { return f.ErrTeapot })

	res, _ := do(t, srv, httptest.NewRequest(http.MethodGet, "/err", nil))

	assert.Equal(t, http.StatusTeapot, res.StatusCode)
}
