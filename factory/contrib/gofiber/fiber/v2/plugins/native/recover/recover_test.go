package recover

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
	srv := newServer(t, NewRecoverWithOptions(&Options{}).Register)
	assert.Equal(t, newServer(t).App().HandlersCount(), srv.App().HandlersCount())
}

func TestRegisterTurnsAPanicIntoAnError(t *testing.T) {
	srv := newServer(t, Register)
	srv.Get("/panic", func(*f.Ctx) error { panic("boom") })

	res, _ := do(t, srv, httptest.NewRequest(http.MethodGet, "/panic", nil))

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
