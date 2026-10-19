package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/examples/multi/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

var google, bing *httptest.Server

// TestMain reads the config of the example and points each client at a
// server of the test that answers with its own status.
func TestMain(m *testing.M) {
	google = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	bing = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	os.Setenv("CONF", "../../configs/config.yaml")
	os.Setenv("BOOST_FACTORY_RESTY_CLIENTS_GOOGLE_HOST", google.URL)
	os.Setenv("BOOST_FACTORY_RESTY_CLIENTS_BING_HOST", bing.URL)
	boost.Start()
	code := m.Run()
	google.Close()
	bing.Close()
	os.Exit(code)
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	for path, want := range map[string]string{
		"/google": `{"url":"` + google.URL + `","status":200}`,
		"/bing":   `{"url":"` + bing.URL + `","status":202}`,
	} {
		rec := httptest.NewRecorder()
		srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		assert.Equal(t, http.StatusOK, rec.Code, path)
		assert.JSONEq(t, want, rec.Body.String(), "each route calls its own client")
	}
}
