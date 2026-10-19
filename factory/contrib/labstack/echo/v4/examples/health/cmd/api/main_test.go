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
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/health/gen/builder/api"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	get := func(target string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
		return rec
	}

	hello := get("/hello")
	assert.Equal(t, http.StatusOK, hello.Code)
	assert.Equal(t, "hello world", hello.Body.String())

	health := get("/health")
	assert.Equal(t, http.StatusOK, health.Code)
	assert.Contains(t, health.Body.String(), `"greetings"`, "the datastore added its checker")
}
