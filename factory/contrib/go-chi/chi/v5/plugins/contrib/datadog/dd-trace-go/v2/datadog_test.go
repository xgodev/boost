package datadog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/config"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// serve builds a server with the plugin from the config defaults and a /x
// route, and answers a GET /x.
func serve(t *testing.T) (*chi.Server, *httptest.ResponseRecorder) {
	t.Helper()
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, Register)
	require.NoError(t, err)
	srv.Get("/x", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	rec := httptest.NewRecorder()
	srv.Mux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	return srv, rec
}

func TestRegisterTracesWhileTheTracerIsEnabled(t *testing.T) {
	srv, rec := serve(t)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Len(t, srv.Mux().Middlewares(), 1)
}

func TestRegisterIsANoOpWhileTheTracerIsDisabled(t *testing.T) {
	// Cleanups run last-in first-out: this reload runs after Setenv restores
	// the variable.
	t.Cleanup(config.Load)
	t.Setenv("BOOST_FACTORY_DATADOG_TRACER_ENABLED", "false")
	config.Load()

	srv, rec := serve(t)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, srv.Mux().Middlewares())
}

func TestDisabledAddsNothing(t *testing.T) {
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, NewDatadogWithOptions(&Options{}).Register)
	require.NoError(t, err)
	assert.Empty(t, srv.Mux().Middlewares())
}
