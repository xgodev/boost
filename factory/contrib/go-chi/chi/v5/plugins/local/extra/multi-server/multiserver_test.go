package multiserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	server "github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// serve builds a server with the plugin from the config defaults and a /x
// route, and answers req.
func serve(t *testing.T, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, Register)
	require.NoError(t, err)
	srv.Get("/x", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	rec := httptest.NewRecorder()
	srv.Mux().ServeHTTP(rec, req)
	return rec
}

func TestDisabledAddsNothing(t *testing.T) {
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, NewMultiServerWithOptions(&Options{}).Register)
	require.NoError(t, err)
	assert.Empty(t, srv.Mux().Middlewares())
	assert.Empty(t, srv.Mux().Routes())
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

	rec := serve(t, httptest.NewRequest(http.MethodGet, "/check", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())

	cancel()
	<-done
	rec = serve(t, httptest.NewRequest(http.MethodGet, "/check", nil))
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}
