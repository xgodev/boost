package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8/examples/connect/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

// TestMain points the client at a server of the test that answers like an
// elasticsearch node.
func TestMain(m *testing.M) {
	es := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.WriteHeader(http.StatusOK)
	}))
	os.Setenv("BOOST_FACTORY_ELASTICSEARCH_ADDRESSES", es.URL)
	boost.Start()
	code := m.Run()
	es.Close()
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

	rec := do(srv, http.MethodGet, "/ping")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"database":"elasticsearch","status":"200 OK"}`, rec.Body.String())

}
