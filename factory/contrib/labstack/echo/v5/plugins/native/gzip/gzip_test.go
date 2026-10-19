package gzip

import (
	gz "compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

var payload = strings.Repeat("boost ", 100)

func newServer(t *testing.T) *echo.Server {
	t.Helper()
	srv, err := echo.NewServerWithOptions(context.Background(), &echo.Options{})
	require.NoError(t, err)
	srv.GET("/x", func(c *e.Context) error { return c.String(http.StatusOK, payload) })
	return srv
}

func TestRegisterWithConfigDefaultsCompressesResponses(t *testing.T) {
	srv := newServer(t)
	require.NoError(t, Register(context.Background(), srv))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set(e.HeaderAcceptEncoding, "gzip")
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "gzip", rec.Header().Get(e.HeaderContentEncoding))
	r, err := gz.NewReader(rec.Body)
	require.NoError(t, err)
	body, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Equal(t, payload, string(body))
}

func TestRegisterRefusesAnInvalidLevel(t *testing.T) {
	err := NewGzipWithOptions(&Options{Enabled: true, Level: 42}).Register(context.Background(), newServer(t))

	assert.Error(t, err)
}
