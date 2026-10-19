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
	"github.com/xgodev/boost/factory/local/extra/middleware/examples/simple/gen/builder/api"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	os.Setenv("CONF", "../../configs/config.yaml")
	boost.Start()
	os.Exit(m.Run())
}

func get(t *testing.T, srv *boostecho.Server, target string) string {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	return rec.Body.String()
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	assert.JSONEq(t, `{"text":"hello, ana"}`, get(t, srv, "/greetings/ana"))
}
