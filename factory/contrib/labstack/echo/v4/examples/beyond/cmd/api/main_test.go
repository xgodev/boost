package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/beyond/gen/builder/api"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/beyond/internal/pkg/domain/service/greeting/greeter"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func TestAPI(t *testing.T) {
	var (
		srv *boostecho.Server
		svc greeter.Service
	)
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv, &svc))...)
	require.NoError(t, app.Err())

	assert.Contains(t, reflect.TypeOf(svc).Elem().PkgPath(), "/gen/builder/middleware/",
		"the handler receives the decorator builder generates for @Middleware(log)")

	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/hello-world", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"Hello World!!"}`, rec.Body.String())
}
