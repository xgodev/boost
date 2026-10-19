package requestid

import (
	"context"
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

func serve(t *testing.T, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	srv, err := echo.NewServerWithOptions(context.Background(), &echo.Options{})
	require.NoError(t, err)
	require.NoError(t, Register(context.Background(), srv))
	srv.GET("/x", func(c *e.Context) error { return c.NoContent(http.StatusOK) })
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, req)
	return rec
}

func TestRegisterWithConfigDefaultsSetsARequestID(t *testing.T) {
	rec := serve(t, httptest.NewRequest(http.MethodGet, "/x", nil))

	assert.NotEmpty(t, rec.Header().Get(e.HeaderXRequestID))
}

func TestRegisterKeepsTheCallersRequestID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set(e.HeaderXRequestID, "abc-123")

	assert.Equal(t, "abc-123", serve(t, req).Header().Get(e.HeaderXRequestID))
}
