package bodylimit

import (
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

func newServer(t *testing.T) *echo.Server {
	t.Helper()
	srv, err := echo.NewServerWithOptions(context.Background(), &echo.Options{})
	require.NoError(t, err)
	srv.POST("/x", func(c *e.Context) error {
		if _, err := io.ReadAll(c.Request().Body); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})
	return srv
}

func post(srv *echo.Server, size int) int {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(strings.Repeat("a", size))))
	return rec.Code
}

func TestOptionsDefaultToEnabledAt8M(t *testing.T) {
	o, err := NewOptions()
	require.NoError(t, err)

	assert.True(t, o.Enabled)
	assert.Equal(t, "8M", o.Size)
}

func TestRegisterRejectsBodiesOverTheLimit(t *testing.T) {
	srv := newServer(t)
	require.NoError(t, NewBodyLimitWithOptions(&Options{Enabled: true, Size: "1K"}).Register(context.Background(), srv))

	assert.Equal(t, http.StatusOK, post(srv, 512))
	assert.Equal(t, http.StatusRequestEntityTooLarge, post(srv, 2048))
}

func TestRegisterRefusesAnUnreadableSize(t *testing.T) {
	err := NewBodyLimitWithOptions(&Options{Enabled: true, Size: "lots"}).Register(context.Background(), newServer(t))

	assert.Error(t, err)
}
