package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache-binary/gen/builder/api"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	os.Setenv("CONF", "../../configs/config.yaml")
	boost.Start()
	os.Exit(m.Run())
}

func do(srv *boostecho.Server, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, req)
	return rec
}

func start(t *testing.T) *boostecho.Server {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())
	return srv
}

func TestAPI(t *testing.T) {
	srv := start(t)

	reading := `{"location":1233,"timestamp":123452123,"temperature":12}`
	rec := do(srv, http.MethodPut, "/readings/1", reading)
	require.Equal(t, http.StatusNoContent, rec.Code, rec.Body.String())

	rec = do(srv, http.MethodGet, "/readings/1", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, reading, rec.Body.String())

	assert.Equal(t, http.StatusNotFound, do(srv, http.MethodGet, "/readings/2", "").Code)
	assert.Equal(t, http.StatusBadRequest, do(srv, http.MethodGet, "/readings/x", "").Code)
}
