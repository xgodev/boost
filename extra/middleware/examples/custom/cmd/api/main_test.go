package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/extra/middleware/examples/custom/gen/builder/api"
	"github.com/xgodev/boost/extra/middleware/examples/custom/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/extra/middleware/examples/custom/internal/pkg/middleware/counter"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"go.uber.org/fx"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func TestAPI(t *testing.T) {
	var srv *boostecho.Server
	var calls *counter.Counter[*entity.Greeting]
	app := fx.New(append(api.Modules(), fx.NopLogger, fx.Populate(&srv), fx.Populate(&calls))...)
	require.NoError(t, app.Err())

	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/greetings/ana", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"text":"hello, ana"}`, rec.Body.String())
	assert.Equal(t, int64(1), calls.Calls(), "the call ran through the counter")
}
