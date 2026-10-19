package echo

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	e "github.com/labstack/echo/v5"
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

// serve builds a server from the config defaults with mutate applied, runs
// Serve in the background and returns its base address once it listens.
// Shutdown runs at cleanup and the test fails if Serve does not return.
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
	srv.GET("/ping", func(c *e.Context) error { return c.String(http.StatusOK, "pong") })

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

func TestOptionsComeFromConfig(t *testing.T) {
	opts, err := NewOptions()
	require.NoError(t, err)

	assert.Equal(t, 8080, opts.Port)
	assert.Equal(t, "REST", opts.Type)
	assert.Equal(t, "HTTP", opts.Protocol)
	assert.True(t, opts.HideBanner)
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
	assert.Equal(t, http.StatusOK, res.StatusCode)
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
	res, body := get(t, client, "https://"+addr+"/ping")
	assert.Equal(t, http.StatusOK, res.StatusCode)
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
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "pong", body)
	assert.Equal(t, 2, res.ProtoMajor, "H2C serves HTTP/2 over cleartext")
}

func TestShutdownBeforeServeMakesServeReturn(t *testing.T) {
	opts, err := NewOptions()
	require.NoError(t, err)
	opts.Port = 0
	srv, err := NewServerWithOptions(context.Background(), opts)
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
	// The wildcard address, the one Serve binds: macOS lets a wildcard bind
	// share a port held on 127.0.0.1 only.
	ln, err := net.Listen("tcp", ":0")
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
		t.Fatal("Serve must return when it cannot listen")
	}
}

func TestServerProxiesRoutesToTheEchoInstance(t *testing.T) {
	srv, err := NewServerWithOptions(context.Background(), &Options{})
	require.NoError(t, err)

	ok := func(c *e.Context) error { return c.NoContent(http.StatusNoContent) }
	srv.GET("/r", ok)
	srv.POST("/r", ok)
	srv.PUT("/r", ok)
	srv.DELETE("/r", ok)
	srv.PATCH("/r", ok)
	srv.OPTIONS("/r", ok)
	srv.HEAD("/r", ok)
	srv.Add(http.MethodGet, "/added", ok)
	srv.Any("/any", ok)
	srv.Match([]string{http.MethodGet}, "/match", ok)
	srv.Group("/g").GET("/r", ok)

	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, "/r"}, {http.MethodPost, "/r"}, {http.MethodPut, "/r"},
		{http.MethodDelete, "/r"}, {http.MethodPatch, "/r"}, {http.MethodOptions, "/r"},
		{http.MethodHead, "/r"}, {http.MethodGet, "/added"}, {http.MethodTrace, "/any"},
		{http.MethodGet, "/match"}, {http.MethodGet, "/g/r"},
	} {
		assert.Equal(t, http.StatusNoContent, statusOf(srv, tc.method, tc.path), "%s %s", tc.method, tc.path)
	}
}

func statusOf(srv *Server, method, path string) int {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec.Code
}
