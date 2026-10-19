package compress

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
	srv := newServer(t, NewCompressWithOptions(&Options{}).Register)
	assert.Equal(t, newServer(t).App().HandlersCount(), srv.App().HandlersCount())
}

func TestRegisterCompressesTheResponse(t *testing.T) {
	srv := newServer(t, Register)
	srv.Get("/big", func(c *f.Ctx) error { return c.SendString(strings.Repeat("a", 4096)) })
	req := httptest.NewRequest(http.MethodGet, "/big", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	res, body := do(t, srv, req)

	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	assert.Less(t, len(body), 4096)
}
