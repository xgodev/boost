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
	"github.com/xgodev/boost/wrapper/config/examples/env_file/gen/builder/api"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	os.Setenv("CONF", "../../configs/config.yaml")
	boost.Start()
	os.Exit(m.Run())
}

func describe(t *testing.T, srv *boostecho.Server) string {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/application", nil))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	return rec.Body.String()
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	// the name from the file, myName from its default
	assert.JSONEq(t, `{"name":"app_test_file","myName":"my_name_test"}`, describe(t, srv))

	// env vars win over the file and the default
	t.Setenv("APP_APPLICATION_NAME", "app_test_env")
	t.Setenv("APP_APPLICATION_MY__NAME", "my_name_test_env")
	config.Load()
	assert.JSONEq(t, `{"name":"app_test_env","myName":"my_name_test_env"}`, describe(t, srv))
}
