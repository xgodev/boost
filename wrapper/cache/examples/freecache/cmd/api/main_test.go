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
	"github.com/xgodev/boost/wrapper/cache/examples/freecache/gen/builder/api"
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

	rec := do(srv, http.MethodPut, "/notes/shopping", `{"text":"milk"}`)
	require.Equal(t, http.StatusNoContent, rec.Code, rec.Body.String())

	rec = do(srv, http.MethodGet, "/notes/shopping", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"key":"shopping","text":"milk"}`, rec.Body.String())

	assert.Equal(t, http.StatusNotFound, do(srv, http.MethodGet, "/notes/todo", "").Code)
}
