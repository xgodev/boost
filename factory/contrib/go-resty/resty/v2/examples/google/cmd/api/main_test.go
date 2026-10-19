package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/examples/google/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

// TestMain reads the config of the example and points the client and its
// health check at a server of the test.
func TestMain(m *testing.M) {
	google := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	os.Setenv("CONF", "../../configs/config.yaml")
	os.Setenv("BOOST_FACTORY_RESTY_CLIENTS_GOOGLE_HOST", google.URL)
	os.Setenv("BOOST_FACTORY_RESTY_CLIENTS_GOOGLE_PLUGINS_HEALTH_HOST", google.URL)
	boost.Start()
	code := m.Run()
	google.Close()
	os.Exit(code)
}

func get(srv *boostecho.Server, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	rec := get(srv, "/google")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"url":"`+os.Getenv("BOOST_FACTORY_RESTY_CLIENTS_GOOGLE_HOST")+`","status":200}`, rec.Body.String())

	rec = get(srv, "/health")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Google Inc", "the resty health plugin registered the checker")
}
