package fiber

import (
	"context"
	stderrors "errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// wait bounds every channel wait in these tests.
const wait = 5 * time.Second

func pong(c *fiber.Ctx) error { return c.SendString("pong") }

// serve builds a server from the config defaults, runs Serve in the
// background and returns its base address once it listens.
func serve(t *testing.T, plugins ...Plugin) (*Server, string) {
	t.Helper()

	opts, err := NewOptions()
	require.NoError(t, err)
	opts.Port = 0

	srv, err := NewServerWithOptions(context.Background(), opts, plugins...)
	require.NoError(t, err)
	srv.Get("/ping", pong)

	done := make(chan struct{})
	go func() {
		srv.Serve(context.Background())
		close(done)
	}()
	t.Cleanup(func() {
		srv.Shutdown(context.Background())
		select {
		case <-done:
		case <-time.After(wait):
			t.Error("Serve did not return after Shutdown")
		}
	})

	select {
	case <-srv.listening:
	case <-done:
		t.Fatal("Serve returned before listening")
	case <-time.After(wait):
		t.Fatal("server did not start listening")
	}
	return srv, srv.addr.String()
}

// do answers a request through the app, without a listener.
func do(t *testing.T, srv *Server, method, path string) (*http.Response, string) {
	t.Helper()
	res, err := srv.App().Test(httptest.NewRequest(method, path, nil), -1)
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return res, string(body)
}

func newServer(t *testing.T, plugins ...Plugin) *Server {
	t.Helper()
	opts, err := NewOptions()
	require.NoError(t, err)
	srv, err := NewServerWithOptions(context.Background(), opts, plugins...)
	require.NoError(t, err)
	return srv
}

func TestOptionsComeFromConfig(t *testing.T) {
	opts, err := NewOptions()
	require.NoError(t, err)

	assert.Equal(t, 8082, opts.Port)
	assert.Equal(t, "REST", opts.Type)
	assert.Equal(t, 4*1024*1024, opts.Config.BodyLimit)
	assert.Equal(t, 256*1024, opts.Config.Concurrency)
	assert.Equal(t, fiber.NetworkTCP4, opts.Config.Network)
	assert.Equal(t, ".fiber.gz", opts.Config.CompressedFileSuffix)
	assert.True(t, opts.Config.DisableStartupMessage)
	assert.Zero(t, opts.Config.ReadTimeout)
}

func TestOptionsWithPathReadTheGivenPath(t *testing.T) {
	t.Cleanup(config.Load)
	t.Setenv("APP_FIBER_PORT", "9090")
	t.Setenv("APP_FIBER_CONFIG_ETAG", "true")
	t.Setenv("APP_FIBER_CONFIG_GET__ONLY", "true")
	t.Setenv("APP_FIBER_CONFIG_READ__TIMEOUT", "3s")

	opts, err := NewOptionsWithPath("app.fiber")
	require.NoError(t, err)

	assert.Equal(t, 9090, opts.Port)
	assert.True(t, opts.Config.ETag)
	assert.True(t, opts.Config.GETOnly)
	assert.Equal(t, 3*time.Second, opts.Config.ReadTimeout)
}

func TestServeServesRoutesUntilShutdown(t *testing.T) {
	srv, addr := serve(t)

	res, err := http.Get("http://" + addr + "/ping")
	require.NoError(t, err)
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "pong", string(body))

	srv.Shutdown(context.Background())
	_, err = http.Get("http://" + addr + "/ping")
	assert.Error(t, err, "the server must not answer after Shutdown")
}

func TestShutdownBeforeServeMakesServeReturn(t *testing.T) {
	srv := newServer(t)

	srv.Shutdown(context.Background())

	done := make(chan struct{})
	go func() {
		srv.Serve(context.Background())
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(wait):
		t.Fatal("Serve after Shutdown must return at once")
	}
}

func TestServeReturnsWhenThePortIsTaken(t *testing.T) {
	ln, err := net.Listen("tcp4", ":0")
	require.NoError(t, err)
	defer ln.Close()

	opts, err := NewOptions()
	require.NoError(t, err)
	opts.Port = ln.Addr().(*net.TCPAddr).Port
	srv, err := NewServerWithOptions(context.Background(), opts)
	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		srv.Serve(context.Background())
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(wait):
		t.Fatal("Serve must return when it cannot bind")
	}
}

func TestPluginsChangeTheConfigBeforeTheAppExists(t *testing.T) {
	srv := newServer(t, func(_ context.Context, s *Server) error {
		assert.Nil(t, s.App(), "the app is created after the plugins")
		s.Config().AppName = "orders"
		return nil
	})

	assert.Equal(t, "orders", srv.App().Config().AppName)
}

func TestPluginMiddlewaresWrapRoutesAddedBeforeThem(t *testing.T) {
	srv := newServer(t,
		func(_ context.Context, s *Server) error {
			s.Get("/x", func(c *fiber.Ctx) error { return c.SendString(c.GetRespHeader("X-Mw")) })
			return nil
		},
		func(_ context.Context, s *Server) error {
			s.Use(func(c *fiber.Ctx) error {
				c.Set("X-Mw", "ran")
				return c.Next()
			})
			return nil
		},
	)

	_, body := do(t, srv, http.MethodGet, "/x")
	assert.Equal(t, "ran", body, "a plugin middleware must wrap a route an earlier plugin added")
}

func TestAFailingPluginStopsTheBuild(t *testing.T) {
	boom := stderrors.New("boom")
	ran := false

	opts, err := NewOptions()
	require.NoError(t, err)
	_, err = NewServerWithOptions(context.Background(), opts,
		func(context.Context, *Server) error { return boom },
		func(context.Context, *Server) error { ran = true; return nil },
	)

	assert.ErrorIs(t, err, boom)
	assert.False(t, ran)
}

func TestRoutesAfterTheBuildAreAddedAtOnce(t *testing.T) {
	srv := newServer(t)
	srv.Post("/p", pong)
	srv.Route("/v1", func(r fiber.Router) { r.Get("/q", pong) })

	res, body := do(t, srv, http.MethodPost, "/p")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "pong", body)
	res, _ = do(t, srv, http.MethodGet, "/v1/q")
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestErrorHandlerJSON(t *testing.T) {
	srv := newServer(t, func(_ context.Context, s *Server) error {
		s.Config().ErrorHandler = ErrorHandlerJSON
		return nil
	})
	srv.All("/notfound", func(*fiber.Ctx) error { return errors.NotFoundf("order 1") })
	srv.All("/fiber", func(*fiber.Ctx) error { return fiber.NewError(http.StatusTeapot, "tea") })
	srv.All("/invalid", func(*fiber.Ctx) error { return validator.ValidationErrors{} })
	ignored := stderrors.New("ignored")
	errors.Ignore(ignored, errors.IgnoreAsSuccess)
	srv.All("/ignored", func(*fiber.Ctx) error { return ignored })

	res, body := do(t, srv, http.MethodGet, "/notfound")
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.JSONEq(t, `{"httpStatusCode":404,"errorCode":"404","message":"order 1 not found"}`, body)

	res, body = do(t, srv, http.MethodGet, "/fiber")
	assert.Equal(t, http.StatusTeapot, res.StatusCode)
	assert.Contains(t, body, `"message":"tea"`)

	res, _ = do(t, srv, http.MethodGet, "/invalid")
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)

	res, body = do(t, srv, http.MethodHead, "/notfound")
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Empty(t, body, "HEAD answers no body")

	res, _ = do(t, srv, http.MethodGet, "/ignored")
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestErrorHandlerString(t *testing.T) {
	srv := newServer(t, func(_ context.Context, s *Server) error {
		s.Config().ErrorHandler = ErrorHandlerString
		return nil
	})
	srv.Get("/x", func(*fiber.Ctx) error { return errors.NotFoundf("order 1") })

	res, body := do(t, srv, http.MethodGet, "/x")
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "order 1 not found", body)
}
