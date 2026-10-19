package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/googlecall/gen/builder/api"
	"go.uber.org/fx"
)

// urlEnv sets internal.pkg.client.page.fetcher.resty.url.
const urlEnv = "INTERNAL_PKG_CLIENT_PAGE_FETCHER_RESTY_URL"

var calls atomic.Int32

// TestMain points the fetcher at a server of the test before boost reads
// the config.
func TestMain(m *testing.M) {
	google := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusAccepted)
	}))
	os.Setenv(urlEnv, google.URL)
	boost.Start()
	code := m.Run()
	google.Close()
	os.Exit(code)
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/google", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"Hello Google!!","status":202}`, rec.Body.String())
	assert.Equal(t, int32(1), calls.Load(), "the handler fetched the page through the resty client")
}
