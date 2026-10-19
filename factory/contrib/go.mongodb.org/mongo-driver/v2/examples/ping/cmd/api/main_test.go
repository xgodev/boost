package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2/examples/ping/gen/builder/api"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func do(srv *boostecho.Server, method, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec
}

// TestAPI needs a live mongo: the client connects when the app starts.
func TestAPI(t *testing.T) {
	if os.Getenv("BOOST_FACTORY_MONGO_URI") == "" {
		t.Skip("set BOOST_FACTORY_MONGO_URI to a mongo to run it, e.g. mongodb://localhost:27017/temp")
	}
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	rec := do(srv, http.MethodGet, "/ping")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"database":"mongo","status":"ok"}`, rec.Body.String())

}
