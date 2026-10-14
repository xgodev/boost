package echo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/model/errors"
	response "github.com/xgodev/boost/model/restresponse"
)

// ErrorHandlerString implements plain text content type error handler.
func ErrorHandlerString(err error, c e.Context) {
	errorHandler(err, c, e.MIMETextPlain)
}

// ErrorHandlerJSON implements JSON content type error handler.
func ErrorHandlerJSON(err error, c e.Context) {
	errorHandler(err, c, e.MIMEApplicationJSON)
}

func errorHandler(err error, c e.Context, contentType string) {
	if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreAsSuccess != 0 {
		if er := c.NoContent(http.StatusOK); er != nil {
			c.Logger().Error(er)
		}
		return
	}

	var (
		status  int
		message string
	)
	if echoErr, ok := err.(*e.HTTPError); ok {
		status = echoErr.Code
		message = fmt.Sprintf("%v", echoErr.Message)
	} else {
		status = ErrorStatusCode(err)
		message = err.Error()
	}

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
		c.Logger().Error(er)
	}
}

// ErrorStatusCode translates err to the respective HTTP status code.
func ErrorStatusCode(err error) int {
	if _, ok := err.(validator.ValidationErrors); ok {
		return http.StatusUnprocessableEntity
	}
	return response.HTTPStatusFor(errors.Classify(err))
}
