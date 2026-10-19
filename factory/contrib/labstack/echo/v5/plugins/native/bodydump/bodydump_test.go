package bodydump

import (
	"context"
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

func TestOptionsDefaultToEnabled(t *testing.T) {
	o, err := NewOptions()
	require.NoError(t, err)

	assert.True(t, o.Enabled)
}

func TestRegisterKeepsTheBodiesFlowing(t *testing.T) {
	srv, err := echo.NewServerWithOptions(context.Background(), &echo.Options{})
	require.NoError(t, err)
	require.NoError(t, Register(context.Background(), srv))
	srv.POST("/x", func(c *e.Context) error {
		var in map[string]string
		if err := c.Bind(&in); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, in)
	})

	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`{"a":"b"}`))
	req.Header.Set(e.HeaderContentType, e.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"a":"b"}`, rec.Body.String())
}
