package prometheus

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

func newServer(t *testing.T, plugins ...chi.Plugin) *chi.Server {
	t.Helper()
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, plugins...)
	require.NoError(t, err)
	srv.Get("/orders/{id}", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	return srv
}

func do(srv *chi.Server, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Mux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestRegisterCountsRequestsByRouteAndServesMetrics(t *testing.T) {
	srv := newServer(t, Register)

	require.Equal(t, http.StatusOK, do(srv, "/orders/7").Code)

	rec := do(srv, "/metrics")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `chi_requests_total{code="200",method="GET",url="/orders/{id}"}`)
	assert.Contains(t, rec.Body.String(), `chi_request_duration_seconds_count{code="200",method="GET",url="/orders/{id}"}`)
}

func TestRegisterOnASecondServerSharesTheCollectors(t *testing.T) {
	newServer(t, Register)

	_, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, Register)

	assert.NoError(t, err, "a second server must not register the chi collectors again")
}

func TestRegisterWithoutTheRouteKeepsTheMiddleware(t *testing.T) {
	o, err := NewOptions()
	require.NoError(t, err)
	o.Route.Enabled = false

	srv := newServer(t, NewPrometheusWithOptions(o).Register)

	assert.Len(t, srv.Mux().Middlewares(), 1)
	assert.Equal(t, http.StatusNotFound, do(srv, "/metrics").Code)
}

func TestDisabledAddsNothing(t *testing.T) {
	srv := newServer(t, NewPrometheusWithOptions(&Options{}).Register)

	assert.Empty(t, srv.Mux().Middlewares())
	assert.Equal(t, http.StatusNotFound, do(srv, "/metrics").Code)
}
