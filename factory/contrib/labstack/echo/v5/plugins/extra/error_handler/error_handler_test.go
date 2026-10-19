package error_handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	e "github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/model/errors"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

func newServer(t *testing.T) *echo.Server {
	t.Helper()
	srv, err := echo.NewServerWithOptions(context.Background(), &echo.Options{Type: "REST"})
	require.NoError(t, err)
	return srv
}

func do(srv *echo.Server, method, target string, body io.Reader) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, httptest.NewRequest(method, target, body))
	return rec
}

func failing(srv *echo.Server, err error) {
	srv.GET("/x", func(c *e.Context) error { return err })
}

func TestRegisterAnswersJSONOnARestServer(t *testing.T) {
	srv := newServer(t)
	require.NoError(t, Register(context.Background(), srv))
	failing(srv, errors.NotFoundf("order 7"))

	rec := do(srv, http.MethodGet, "/x", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"httpStatusCode":404,"errorCode":"404","message":"order 7 not found"}`, rec.Body.String())
}

func TestRegisterAnswersPlainTextOnAnotherServerType(t *testing.T) {
	srv, err := echo.NewServerWithOptions(context.Background(), &echo.Options{Type: "GRPC"})
	require.NoError(t, err)
	require.NoError(t, Register(context.Background(), srv))
	failing(srv, errors.NotFoundf("order 7"))

	rec := do(srv, http.MethodGet, "/x", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "order 7 not found", rec.Body.String())
}

func TestRegisterLeavesACommittedResponseAlone(t *testing.T) {
	srv := newServer(t)
	require.NoError(t, Register(context.Background(), srv))
	srv.GET("/x", func(c *e.Context) error {
		if err := c.String(http.StatusAccepted, "partial"); err != nil {
			return err
		}
		return errors.Internalf("late failure")
	})

	rec := do(srv, http.MethodGet, "/x", nil)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Equal(t, "partial", rec.Body.String())
}
