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
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/multiecho/gen/builder/api"
	"go.uber.org/fx"
)

// TestMain lets both servers listen on a free port: the test starts the app.
func TestMain(m *testing.M) {
	os.Setenv("BOOST_FACTORY_ECHO_PORT", "0")
	os.Setenv("API_SERVER_SECOND_PORT", "0")
	boost.Start()
	os.Exit(m.Run())
}

func get(srv *boostecho.Server, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestAPI(t *testing.T) {
	var rest, second *boostecho.Server
	app := fx.New(append(api.Modules(),
		fx.NopLogger,
		fx.Populate(&rest),
		fx.Invoke(fx.Annotate(func(srv *boostecho.Server) { second = srv }, fx.ParamTags(`name:"second"`))),
	)...)
	require.NoError(t, app.Err())
	require.NotSame(t, rest, second, "each server is its own echo instance")

	require.NoError(t, app.Start(context.Background()))
	t.Cleanup(func() { _ = app.Stop(context.Background()) })
	require.Eventually(t, func() bool {
		return rest.Instance().ListenerAddr() != nil && second.Instance().ListenerAddr() != nil
	}, 5*time.Second, 10*time.Millisecond, "the app serves both servers")

	for name, srv := range map[string]*boostecho.Server{"rest": rest, "second": second} {
		rec := get(srv, "/test")
		assert.Equal(t, http.StatusOK, rec.Code, name)
		assert.JSONEq(t, `{"message":"Hello World!!"}`, rec.Body.String(), name)
	}

	assert.Equal(t, http.StatusOK, get(second, "/check").Code, "both servers are up")
	assert.Equal(t, http.StatusNotFound, get(rest, "/check").Code, "the multi-server plugin is on second only")
}
