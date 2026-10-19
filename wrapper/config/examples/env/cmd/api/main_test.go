package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/config/examples/env/gen/builder/api"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func name(t *testing.T, srv *boostecho.Server) string {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/application", nil))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	return rec.Body.String()
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	// the default of config.Add
	assert.JSONEq(t, `{"name":"app_test"}`, name(t, srv))

	// an env var, once config.Load reads it again
	t.Setenv("APP_APPLICATION_NAME", "app_test_env_2")
	config.Load()
	assert.JSONEq(t, `{"name":"app_test_env_2"}`, name(t, srv))

	t.Setenv("APP_APPLICATION_NAME", "app_test_env_3")
	config.Load()
	assert.JSONEq(t, `{"name":"app_test_env_3"}`, name(t, srv))
}
