package log

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
	"github.com/xgodev/boost/wrapper/log"
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

func TestRegisterHandlesTheErrorOnceAndLogsTheStatus(t *testing.T) {
	srv := newServer(t)
	handled := 0
	srv.Instance().HTTPErrorHandler = func(c *e.Context, err error) {
		handled++
		_ = c.String(e.StatusCode(err), err.Error())
	}
	require.NoError(t, Register(context.Background(), srv))
	srv.GET("/x", func(c *e.Context) error { return e.ErrNotFound })

	rec := do(srv, http.MethodGet, "/x", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, 1, handled, "the middleware handles the error and returns nil")
}

func TestRegisterPutsTheRequestLoggerInTheContext(t *testing.T) {
	srv := newServer(t)
	require.NoError(t, Register(context.Background(), srv))
	var got log.Logger
	srv.GET("/x", func(c *e.Context) error {
		got = log.FromContext(c.Request().Context())
		return c.NoContent(http.StatusOK)
	})

	require.Equal(t, http.StatusOK, do(srv, http.MethodGet, "/x", nil).Code)
	assert.NotNil(t, got)
}
