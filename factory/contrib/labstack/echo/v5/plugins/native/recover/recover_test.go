package recover

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

func TestRegisterWithConfigDefaultsTurnsPanicsInto500(t *testing.T) {
	srv, err := echo.NewServerWithOptions(context.Background(), &echo.Options{})
	require.NoError(t, err)
	require.NoError(t, Register(context.Background(), srv))
	srv.GET("/x", func(*e.Context) error { panic("boom") })

	rec := httptest.NewRecorder()
	assert.NotPanics(t, func() {
		srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	})
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
