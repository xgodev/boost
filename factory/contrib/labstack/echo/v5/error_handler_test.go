package echo

import (
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-playground/validator/v10"
	e "github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/model/errors"
	response "github.com/xgodev/boost/model/restresponse"
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

	errorHandler(c, sentinel, e.MIMEApplicationJSON)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestErrorHandlerJSONWritesTheStatusAndEnvelope(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		status  int
		message string
	}{
		{"echo HTTPError", e.NewHTTPError(http.StatusTeapot, "short and stout"), http.StatusTeapot, "short and stout"},
		{"echo HTTPError without message", &e.HTTPError{Code: http.StatusConflict}, http.StatusConflict, http.StatusText(http.StatusConflict)},
		{"echo sentinel", e.ErrNotFound, http.StatusNotFound, http.StatusText(http.StatusNotFound)},
		{"boost kind", errors.NotFoundf("order 7"), http.StatusNotFound, "order 7 not found"},
		{"validation", validator.ValidationErrors{}, http.StatusUnprocessableEntity, ""},
		{"plain", stderrors.New("boom"), http.StatusInternalServerError, "boom"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := e.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

			ErrorHandlerJSON(c, tc.err)

			require.Equal(t, tc.status, rec.Code)
			var body response.Error
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tc.status, body.HttpStatusCode)
			assert.Equal(t, strconv.Itoa(tc.status), body.ErrorCode)
			if tc.message != "" {
				assert.Equal(t, tc.message, body.Message)
			}
		})
	}
}

func TestErrorHandlerStringWritesPlainText(t *testing.T) {
	rec := httptest.NewRecorder()
	c := e.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

	ErrorHandlerString(c, e.NewHTTPError(http.StatusBadRequest, "bad"))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "bad", rec.Body.String())
	assert.Contains(t, rec.Header().Get(e.HeaderContentType), e.MIMETextPlain)
}

func TestErrorHandlerWritesNoBodyOnHead(t *testing.T) {
	rec := httptest.NewRecorder()
	c := e.New().NewContext(httptest.NewRequest(http.MethodHead, "/", nil), rec)

	ErrorHandlerJSON(c, errors.NotFoundf("x"))

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Empty(t, rec.Body.String())
}
