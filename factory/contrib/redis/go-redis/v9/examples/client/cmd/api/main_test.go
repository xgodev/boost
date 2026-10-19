package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9/examples/client/gen/builder/api"
	"go.uber.org/fx"
)

var redis *miniredis.Miniredis

// TestMain points the client at an in-memory redis.
func TestMain(m *testing.M) {
	redis = miniredis.NewMiniRedis()
	if err := redis.Start(); err != nil {
		panic(err)
	}
	os.Setenv("BOOST_FACTORY_REDIS_CLIENT_ADDR", redis.Addr())
	boost.Start()
	code := m.Run()
	redis.Close()
	os.Exit(code)
}

func do(srv *boostecho.Server, method, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	for want := 1; want <= 2; want++ {
		rec := do(srv, http.MethodPost, "/counters/visits")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"name":"visits","count":%d}`, want), rec.Body.String())
	}
	got, err := redis.Get("counter:visits")
	require.NoError(t, err)
	assert.Equal(t, "2", got, "the count is kept on redis")

	rec := do(srv, http.MethodGet, "/health")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"redis"`, "the redis health plugin registered the checker")
}
