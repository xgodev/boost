package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/cache/examples/service/gen/builder/api"
	"github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/cache/examples/service/internal/app/api/rest/handler/report/finder"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	os.Setenv("CONF", "../../configs/config.yaml")
	boost.Start()
	os.Exit(m.Run())
}

func get(srv *boostecho.Server, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func report(t *testing.T, srv *boostecho.Server, id string) finder.Response {
	rec := get(srv, "/reports/"+id)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var r finder.Response
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &r))
	return r
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv))...)
	require.NoError(t, app.Err())

	first := report(t, srv, "1")
	assert.Equal(t, "1", first.ID)
	// the second call answers from the cache: same build
	assert.Equal(t, first, report(t, srv, "1"))
	// another id is another build
	assert.NotEqual(t, first.Build, report(t, srv, "2").Build)

	// a failed build goes through the fallback and is not cached
	assert.Equal(t, http.StatusBadRequest, get(srv, "/reports/abc").Code)
	assert.Equal(t, http.StatusBadRequest, get(srv, "/reports/abc").Code)
}
