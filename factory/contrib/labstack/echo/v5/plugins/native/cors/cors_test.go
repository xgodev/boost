package cors

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

func preflight(t *testing.T, srv *echo.Server) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodOptions, "/x", nil)
	req.Header.Set(e.HeaderOrigin, "http://a.test")
	req.Header.Set(e.HeaderAccessControlRequestMethod, http.MethodGet)
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, req)
	return rec
}

func newServer(t *testing.T) *echo.Server {
	t.Helper()
	srv, err := echo.NewServerWithOptions(context.Background(), &echo.Options{})
	require.NoError(t, err)
	srv.GET("/x", func(c *e.Context) error { return c.NoContent(http.StatusOK) })
	return srv
}

func TestRegisterWithConfigDefaultsAllowsAnyOrigin(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, Register(context.Background(), srv))

	rec := preflight(t, srv)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "*", rec.Header().Get(e.HeaderAccessControlAllowOrigin))
	assert.Equal(t, "5200", rec.Header().Get(e.HeaderAccessControlMaxAge))
}

func TestRegisterRefusesTheStarOriginWithCredentials(t *testing.T) {
	o, err := NewOptions()
	require.NoError(t, err)
	o.Allowed.Credentials = true

	err = NewCorsWithOptions(o).Register(context.Background(), newServer(t))

	assert.Error(t, err)
}

func TestRegisterDisabledAddsNoHeaders(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, NewCorsWithOptions(&Options{}).Register(context.Background(), srv))

	assert.Empty(t, preflight(t, srv).Header().Get(e.HeaderAccessControlAllowOrigin))
}
