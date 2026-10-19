package datadog

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

func TestRegisterIsANoOpWhileTheTracerIsDisabled(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, Register(context.Background(), srv))

	srv.GET("/x", func(c *e.Context) error { return c.NoContent(http.StatusOK) })
	assert.Equal(t, http.StatusOK, do(srv, http.MethodGet, "/x", nil).Code)
}
