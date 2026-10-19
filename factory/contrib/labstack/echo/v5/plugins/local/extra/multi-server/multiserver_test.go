package multiserver

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	server "github.com/xgodev/boost/extra/multiserver"
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
	srv := newServer(t)
	require.NoError(t, Register(context.Background(), srv))

	ctx, cancel := context.WithCancel(context.Background())
	s := stub{started: make(chan struct{})}
	done := make(chan struct{})
	go func() {
		server.Serve(ctx, s)
		close(done)
	}()
	<-s.started

	rec := do(srv, http.MethodGet, "/check", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())

	cancel()
	<-done

	rec = do(srv, http.MethodGet, "/check", nil)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Equal(t, "Service Unavailable", rec.Body.String())
}

func TestRegisterDisabledAddsNoRoute(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, NewMultiServerWithOptions(&Options{Route: "/check"}).Register(context.Background(), srv))

	assert.Equal(t, http.StatusNotFound, do(srv, http.MethodGet, "/check", nil).Code)
}
