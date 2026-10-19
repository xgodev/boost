package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/multiserver/gen/builder/api"
	"go.uber.org/fx"
)

// TestMain lets echo listen on a free port and makes the background server stop
// soon: the test starts the app.
func TestMain(m *testing.M) {
	os.Setenv("BOOST_FACTORY_ECHO_PORT", "0")
	os.Setenv("INTERNAL_APP_API_BACKGROUND_DURATION", "300ms")
	boost.Start()
	os.Exit(m.Run())
}

func get(srv *boostecho.Server, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestAPI(t *testing.T) {
	var rest *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&rest))...)
	require.NoError(t, app.Err())

	require.NoError(t, app.Start(context.Background()))
	t.Cleanup(func() { _ = app.Stop(context.Background()) })
	require.Eventually(t, func() bool {
		return rest.Instance().ListenerAddr() != nil
	}, 5*time.Second, 10*time.Millisecond, "the app serves echo")

	rec := get(rest, "/test")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"Hello World!!"}`, rec.Body.String())

	assert.Equal(t, http.StatusOK, get(rest, "/check").Code, "echo and the background server are up")
	assert.Eventually(t, func() bool {
		return get(rest, "/check").Code == http.StatusServiceUnavailable
	}, 5*time.Second, 20*time.Millisecond, "the background server stopped")
	assert.Equal(t, http.StatusOK, get(rest, "/test").Code, "echo keeps serving")
}
