package health

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

func TestRegisterServesTheHealthRoute(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, Register(context.Background(), srv))

	rec := do(srv, http.MethodGet, "/health", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"OK"`)
}

func TestRegisterUsesTheConfiguredRoute(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, NewHealthWithOptions(&Options{Enabled: true, Route: "/live"}).Register(context.Background(), srv))

	assert.Equal(t, http.StatusOK, do(srv, http.MethodGet, "/live", nil).Code)
	assert.Equal(t, http.StatusNotFound, do(srv, http.MethodGet, "/health", nil).Code)
}

func TestRegisterDisabledAddsNoRoute(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, NewHealthWithOptions(&Options{Route: "/health"}).Register(context.Background(), srv))

	assert.Equal(t, http.StatusNotFound, do(srv, http.MethodGet, "/health", nil).Code)
}
