package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/examples/boot/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	os.Setenv("CONF", "../../configs/config.yaml")
	boost.Start()
	os.Exit(m.Run())
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/greetings/world", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"text":"hello, world"}`, rec.Body.String())
}
