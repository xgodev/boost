package echo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/model/errors"
)

type customNotFound struct{ msg string }

func (c *customNotFound) Error() string { return c.msg }

func TestErrorStatusCode_RegisteredCustom(t *testing.T) {
	defer errors.ResetRegistry()
	errors.RegisterMatch(func(err error) bool {
		_, ok := err.(*customNotFound)
		return ok
	}, errors.KindNotFound)

	if got := ErrorStatusCode(&customNotFound{"missing"}); got != http.StatusNotFound {
		t.Fatalf("ErrorStatusCode = %d, want 404", got)
	}
}

func TestErrorStatusCode_BuiltinNotFound(t *testing.T) {
	if got := ErrorStatusCode(errors.NotFoundf("x")); got != http.StatusNotFound {
		t.Fatalf("ErrorStatusCode = %d, want 404", got)
	}
}

func TestErrorHandler_IgnoreAsSuccess(t *testing.T) {
	defer errors.ResetRegistry()
	sentinel := errors.New("ignorable")
	errors.Ignore(sentinel, errors.IgnoreAsSuccess)

	rec := httptest.NewRecorder()
	srv := e.New()
	c := srv.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

	errorHandler(sentinel, c, e.MIMEApplicationJSON)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
