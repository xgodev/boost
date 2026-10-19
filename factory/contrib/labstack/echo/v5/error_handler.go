package echo

import (
	stderrors "errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	e "github.com/labstack/echo/v5"
	"github.com/xgodev/boost/model/errors"
	response "github.com/xgodev/boost/model/restresponse"
)

// ErrorHandlerString implements plain text content type error handler.
func ErrorHandlerString(c *e.Context, err error) {
	errorHandler(c, err, e.MIMETextPlain)
}

// ErrorHandlerJSON implements JSON content type error handler.
func ErrorHandlerJSON(c *e.Context, err error) {
	errorHandler(c, err, e.MIMEApplicationJSON)
}

func errorHandler(c *e.Context, err error, contentType string) {
	if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreAsSuccess != 0 {
		if er := c.NoContent(http.StatusOK); er != nil {
			c.Logger().Error(er.Error())
		}
		return
	}

	status, message := statusAndMessage(err)

	var er error
	if c.Request().Method == http.MethodHead {
		er = c.NoContent(status)
	} else {
		switch contentType {
		case e.MIMEApplicationJSON:
			er = c.JSON(status, response.Error{HttpStatusCode: status, ErrorCode: strconv.Itoa(status), Message: message})
		default:
			er = c.String(status, message)
		}
	}
	if er != nil {
		c.Logger().Error(er.Error())
	}
}

// statusAndMessage reads the status echo attached to err (an *echo.HTTPError,
// or a sentinel like echo.ErrNotFound) and falls back to ErrorStatusCode.
func statusAndMessage(err error) (int, string) {
	var he *e.HTTPError
	if stderrors.As(err, &he) {
		if he.Message == "" {
			return he.Code, http.StatusText(he.Code)
		}
		return he.Code, he.Message
	}
	if status := e.StatusCode(err); status != 0 {
		return status, err.Error()
	}
	return ErrorStatusCode(err), err.Error()
}

// ErrorStatusCode translates err to the respective HTTP status code.
func ErrorStatusCode(err error) int {
	if _, ok := err.(validator.ValidationErrors); ok {
		return http.StatusUnprocessableEntity
	}
	return response.HTTPStatusFor(errors.Classify(err))
}
