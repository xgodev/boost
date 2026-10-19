package otel

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	e "github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
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

func TestRegisterTracesRequestsWithTheBoostProvider(t *testing.T) {
	spans := tracetest.NewSpanRecorder()
	previous := otelboost.TracerProvider
	otelboost.TracerProvider = sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spans))
	t.Cleanup(func() { otelboost.TracerProvider = previous })

	srv := newServer(t)
	require.NoError(t, Register(context.Background(), srv))
	srv.GET("/orders/:id", func(c *e.Context) error { return c.NoContent(http.StatusOK) })

	require.Equal(t, http.StatusOK, do(srv, http.MethodGet, "/orders/7", nil).Code)

	ended := spans.Ended()
	require.Len(t, ended, 1)
	assert.Equal(t, "GET /orders/:id", ended[0].Name())
}

func TestRegisterDisabledAddsNoMiddleware(t *testing.T) {
	spans := tracetest.NewSpanRecorder()
	previous := otelboost.TracerProvider
	otelboost.TracerProvider = sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spans))
	t.Cleanup(func() { otelboost.TracerProvider = previous })

	srv := newServer(t)
	require.NoError(t, NewOtelWithOptions(&Options{}).Register(context.Background(), srv))
	srv.GET("/x", func(c *e.Context) error { return c.NoContent(http.StatusOK) })

	require.Equal(t, http.StatusOK, do(srv, http.MethodGet, "/x", nil).Code)
	assert.Empty(t, spans.Ended())
}
