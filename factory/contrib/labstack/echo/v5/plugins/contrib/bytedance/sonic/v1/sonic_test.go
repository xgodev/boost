package json

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	e "github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
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

type order struct {
	ID int `json:"id"`
}

func echoOrder(srv *echo.Server) {
	srv.POST("/orders", func(c *e.Context) error {
		var o order
		if err := c.Bind(&o); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, o)
	})
}

func post(srv *echo.Server, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	req.Header.Set(e.HeaderContentType, e.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	srv.Instance().ServeHTTP(rec, req)
	return rec
}

func TestRegisterSerializesWithThisEncoder(t *testing.T) {
	srv := newServer(t)
	require.NoError(t, Register(context.Background(), srv))
	echoOrder(srv)

	rec := post(srv, `{"id":7}`)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"id":7}`, rec.Body.String())
	assert.IsType(t, JSONSerializer{}, srv.Instance().JSONSerializer)
}

func TestDeserializeAnswersBadRequestOnMalformedJSON(t *testing.T) {
	srv := newServer(t)
	require.NoError(t, Register(context.Background(), srv))
	echoOrder(srv)

	assert.Equal(t, http.StatusBadRequest, post(srv, `{"id":`).Code)
	assert.Equal(t, http.StatusBadRequest, post(srv, `{"id":"seven"}`).Code)
}

func TestRegisterDisabledKeepsEchoEncoder(t *testing.T) {
	srv := newServer(t)

	require.NoError(t, NewSonicWithOptions(&Options{}).Register(context.Background(), srv))

	assert.NotEqual(t, JSONSerializer{}, srv.Instance().JSONSerializer)
}
