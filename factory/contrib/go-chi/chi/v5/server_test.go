package chi

import (
	"context"
	"crypto/tls"
	stderrors "errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	c "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"golang.org/x/net/http2"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// wait bounds every channel wait in these tests.
const wait = 5 * time.Second

func pong(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("pong")) }

// serve builds a server from the config defaults with mutate applied, runs
// Serve in the background and returns its base address once it listens.
func serve(t *testing.T, mutate func(*Options), plugins ...Plugin) (*Server, string) {
	t.Helper()

	opts, err := NewOptions()
	require.NoError(t, err)
	opts.Port = 0
	if mutate != nil {
		mutate(opts)
	}

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

func get(t *testing.T, client *http.Client, url string) (*http.Response, string) {
	t.Helper()
	res, err := client.Get(url)
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return res, string(body)
}

func do(srv *Server, method, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Mux().ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func TestOptionsComeFromConfig(t *testing.T) {
	opts, err := NewOptions()
	require.NoError(t, err)

	assert.Equal(t, 8080, opts.Port)
	assert.Equal(t, "HTTP", opts.Protocol)
	assert.False(t, opts.TLS.Enabled)
	assert.Equal(t, "AUTO", opts.TLS.Type)
	assert.Equal(t, "localhost", opts.TLS.Auto.Host)
}

func TestServeServesRoutesUntilShutdown(t *testing.T) {
	srv, addr := serve(t, nil)

	res, body := get(t, http.DefaultClient, "http://"+addr+"/ping")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "pong", body)

	srv.Shutdown(context.Background())
	_, err := http.Get("http://" + addr + "/ping")
	assert.Error(t, err, "the server must not answer after Shutdown")
}

func TestServeServesTLSWithAnAutoCertificate(t *testing.T) {
	_, addr := serve(t, func(o *Options) { o.TLS.Enabled = true })

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // self-signed AUTO certificate
		ForceAttemptHTTP2: true,
	}}
	res, body := get(t, client, "https://"+addr+"/ping")
	assert.Equal(t, "pong", body)
	assert.Equal(t, 2, res.ProtoMajor, "TLS serves HTTP/2 unless disableHTTP2")
}

func TestServeTLSDisableHTTP2ServesHTTP1(t *testing.T) {
	_, addr := serve(t, func(o *Options) { o.TLS.Enabled = true; o.DisableHTTP2 = true })

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // self-signed AUTO certificate
		ForceAttemptHTTP2: true,
	}}
	res, _ := get(t, client, "https://"+addr+"/ping")
	assert.Equal(t, 1, res.ProtoMajor)
}

func TestServeServesTLSFromFiles(t *testing.T) {
	cert, key, err := generateCertificate("localhost")
	require.NoError(t, err)
	dir := t.TempDir()
	certFile, keyFile := filepath.Join(dir, "cert.pem"), filepath.Join(dir, "key.pem")
	require.NoError(t, os.WriteFile(certFile, cert, 0o600))
	require.NoError(t, os.WriteFile(keyFile, key, 0o600))

	_, addr := serve(t, func(o *Options) {
		o.TLS.Enabled = true
		o.TLS.Type = "FILE"
		o.TLS.File.Cert = certFile
		o.TLS.File.Key = keyFile
	})

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // self-signed test certificate
	}}
	_, body := get(t, client, "https://"+addr+"/ping")
	assert.Equal(t, "pong", body)
}

func TestServeServesH2C(t *testing.T) {
	_, addr := serve(t, func(o *Options) { o.Protocol = "H2C" })

	client := &http.Client{Transport: &http2.Transport{
		AllowHTTP: true,
		DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, network, addr)
		},
	}}
	res, body := get(t, client, "http://"+addr+"/ping")
	assert.Equal(t, "pong", body)
	assert.Equal(t, 2, res.ProtoMajor, "H2C serves HTTP/2 over cleartext")
}

func TestShutdownBeforeServeMakesServeReturn(t *testing.T) {
	srv, err := NewServerWithOptions(context.Background(), &Options{})
	require.NoError(t, err)

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
	ln, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer ln.Close()

	srv, err := NewServerWithOptions(context.Background(), &Options{Port: ln.Addr().(*net.TCPAddr).Port})
	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		srv.Serve(context.Background())
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(wait):
		t.Fatal("Serve must return when it cannot listen")
	}
}

func header(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Mw", name)
			next.ServeHTTP(w, r)
		})
	}
}

func TestPluginMiddlewaresWrapTheRoutesOfEarlierPlugins(t *testing.T) {
	route := func(_ context.Context, s *Server) error { s.Get("/health", pong); return nil }
	mw := func(_ context.Context, s *Server) error { s.Use(header("cors")); return nil }

	srv, err := NewServerWithOptions(context.Background(), &Options{}, route, mw)
	require.NoError(t, err, "chi panics on a middleware after a route; the server must order them")

	rec := do(srv, http.MethodGet, "/health")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "cors", rec.Header().Get("X-Mw"))

	srv.Get("/app", pong)
	assert.Equal(t, "cors", do(srv, http.MethodGet, "/app").Header().Get("X-Mw"), "routes of the app run the plugin middlewares too")
}

func TestPluginErrorStopsTheBuild(t *testing.T) {
	boom := stderrors.New("boom")
	ran := false
	_, err := NewServerWithOptions(context.Background(), &Options{},
		func(context.Context, *Server) error { return boom },
		func(context.Context, *Server) error { ran = true; return nil },
	)
	assert.ErrorIs(t, err, boom)
	assert.False(t, ran)
}

func TestServerProxiesRoutesToTheMux(t *testing.T) {
	srv, err := NewServerWithOptions(context.Background(), &Options{})
	require.NoError(t, err)

	ok := func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }
	srv.Get("/r", ok)
	srv.Post("/r", ok)
	srv.Put("/r", ok)
	srv.Delete("/r", ok)
	srv.Patch("/r", ok)
	srv.Options("/r", ok)
	srv.Head("/r", ok)
	srv.Connect("/r", ok)
	srv.Trace("/r", ok)
	srv.MethodFunc(http.MethodGet, "/method-func", ok)
	srv.Method(http.MethodGet, "/method", http.HandlerFunc(ok))
	srv.HandleFunc("/handle-func", ok)
	srv.Handle("/handle", http.HandlerFunc(ok))
	srv.Route("/route", func(r c.Router) { r.Get("/r", ok) })
	srv.Group(func(r c.Router) { r.Get("/group", ok) })
	srv.Mount("/mount", http.HandlerFunc(ok))
	srv.With(header("with")).Get("/with", ok)

	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, "/r"}, {http.MethodPost, "/r"}, {http.MethodPut, "/r"},
		{http.MethodDelete, "/r"}, {http.MethodPatch, "/r"}, {http.MethodOptions, "/r"},
		{http.MethodHead, "/r"}, {http.MethodConnect, "/r"}, {http.MethodTrace, "/r"},
		{http.MethodGet, "/method-func"}, {http.MethodGet, "/method"},
		{http.MethodPost, "/handle-func"}, {http.MethodPut, "/handle"},
		{http.MethodGet, "/route/r"}, {http.MethodGet, "/group"}, {http.MethodGet, "/mount"},
		{http.MethodGet, "/with"},
	} {
		assert.Equal(t, http.StatusNoContent, do(srv, tc.method, tc.path).Code, "%s %s", tc.method, tc.path)
	}
}

func TestServerWithConfigPathReadsItsOwnKeys(t *testing.T) {
	t.Setenv("APP_ADMIN_PORT", "9191")
	srv, err := NewServerWithConfigPath(context.Background(), "app.admin")
	require.NoError(t, err)
	assert.Equal(t, 9191, srv.ServerOptions().Port)
	assert.Equal(t, "HTTP", srv.ServerOptions().Protocol)
}
