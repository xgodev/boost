package main

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/extra/health/examples/simple/gen/builder/api"
	"github.com/xgodev/boost/extra/health/examples/simple/internal/app/api/rest/handler/health/checker"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

var dependency net.Listener

// TestMain puts the dependency on a free port.
func TestMain(m *testing.M) {
	var err error
	dependency, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	os.Setenv("APP_DEPENDENCY_ADDR", dependency.Addr().String())
	boost.Start()
	os.Exit(m.Run())
}

func health(t *testing.T, srv *boostecho.Server) (int, []checker.Check) {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	var checks []checker.Check
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &checks))
	return rec.Code, checks
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	code, checks := health(t, srv)
	assert.Equal(t, http.StatusOK, code)
	require.Len(t, checks, 1)
	assert.Equal(t, checker.Check{Name: "dependency", OK: true}, checks[0])

	require.NoError(t, dependency.Close())
	code, checks = health(t, srv)
	assert.Equal(t, http.StatusServiceUnavailable, code, "the required dependency is down")
	require.Len(t, checks, 1)
	assert.False(t, checks[0].OK)
	assert.NotEmpty(t, checks[0].Error)
}
