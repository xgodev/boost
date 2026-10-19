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
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/notfound/gen/builder/api"
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

	hello := get("/items/hello")
	assert.Equal(t, http.StatusOK, hello.Code)
	assert.Equal(t, "hello world", hello.Body.String())

	missing := get("/items/missing")
	assert.Equal(t, http.StatusNotFound, missing.Code, "error_handler answers a NotFound with 404")
	assert.Contains(t, missing.Body.String(), "item missing")

	metrics := get("/metrics")
	assert.Equal(t, http.StatusOK, metrics.Code, "the prometheus plugin serves its route")
	assert.Contains(t, metrics.Body.String(), "/items/:id", "the requests above were measured by route")
}
