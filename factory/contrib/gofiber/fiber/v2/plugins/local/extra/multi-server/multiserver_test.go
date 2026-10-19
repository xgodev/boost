package multiserver

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
	server "github.com/xgodev/boost/extra/multiserver"
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
	srv := newServer(t, NewMultiServerWithOptions(&Options{}).Register)
	assert.Equal(t, newServer(t).App().HandlersCount(), srv.App().HandlersCount())
}

// stub is a multiserver server that runs until its context is done.
type stub struct{ started chan struct{} }

func (s stub) Serve(ctx context.Context) {
	close(s.started)
	<-ctx.Done()
}

func (stub) Shutdown(context.Context) {}

// ran: multiserver keeps every served server for the process, so a rerun
// (-count>1) would still see this test's server down.
var ran bool

func TestRegisterServesTheMultiserverState(t *testing.T) {
	if ran {
		t.Skip("multiserver state is process-wide; the first run covers it")
	}
	ran = true

	ctx, cancel := context.WithCancel(context.Background())
	s := stub{started: make(chan struct{})}
	done := make(chan struct{})
	go func() {
		server.Serve(ctx, s)
		close(done)
	}()
	<-s.started

	res, body := do(t, newServer(t, Register), httptest.NewRequest(http.MethodGet, "/check", nil))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "OK", body)

	cancel()
	<-done

	res, _ = do(t, newServer(t, Register), httptest.NewRequest(http.MethodGet, "/check", nil))
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}
