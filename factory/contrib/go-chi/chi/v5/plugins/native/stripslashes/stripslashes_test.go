package stripslashes

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
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, NewStripSlashesWithOptions(&Options{}).Register)
	require.NoError(t, err)
	assert.Empty(t, srv.Mux().Middlewares())
	assert.Empty(t, srv.Mux().Routes())
}

func TestRegisterRoutesATrailingSlash(t *testing.T) {
	assert.Equal(t, http.StatusOK, serve(t, httptest.NewRequest(http.MethodGet, "/x/", nil)).Code)
}
