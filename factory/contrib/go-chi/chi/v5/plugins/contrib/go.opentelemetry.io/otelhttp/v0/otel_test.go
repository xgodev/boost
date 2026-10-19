package otel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	c "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// record swaps the boost tracer provider for one that keeps the ended spans.
func record(t *testing.T) *tracetest.SpanRecorder {
	t.Helper()
	spans := tracetest.NewSpanRecorder()
	previous := otelboost.TracerProvider
	otelboost.TracerProvider = sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spans))
	t.Cleanup(func() { otelboost.TracerProvider = previous })
	return spans
}

func do(srv *chi.Server, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Mux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestRegisterTracesRequestsByRouteWithTheBoostProvider(t *testing.T) {
	spans := record(t)
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, Register)
	require.NoError(t, err)
	srv.Get("/orders/{id}", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	require.Equal(t, http.StatusOK, do(srv, "/orders/7").Code)

	ended := spans.Ended()
	require.Len(t, ended, 1)
	assert.Equal(t, "GET /orders/{id}", ended[0].Name())
}

func TestRegisterNamesSpansOfMountedRouters(t *testing.T) {
	spans := record(t)
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, Register)
	require.NoError(t, err)
	srv.Route("/api", func(r c.Router) {
		r.Get("/users/{id}", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	})

	require.Equal(t, http.StatusOK, do(srv, "/api/users/3").Code)

	ended := spans.Ended()
	require.Len(t, ended, 1)
	assert.Equal(t, "GET /api/users/{id}", ended[0].Name())
}

func TestDisabledAddsNothing(t *testing.T) {
	spans := record(t)
	srv, err := chi.NewServerWithOptions(context.Background(), &chi.Options{}, NewOtelWithOptions(&Options{}).Register)
	require.NoError(t, err)
	srv.Get("/x", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	require.Equal(t, http.StatusOK, do(srv, "/x").Code)
	assert.Empty(t, srv.Mux().Middlewares())
	assert.Empty(t, spans.Ended())
}
