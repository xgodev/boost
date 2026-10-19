package fiber

import (
	stderrors "errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/model/errors"
	response "github.com/xgodev/boost/model/restresponse"
)

// ErrorHandlerString implements plain text content type error handler.
func ErrorHandlerString(c *fiber.Ctx, err error) error {
	return errorHandler(c, err, false)
}

// ErrorHandlerJSON implements JSON content type error handler.
func ErrorHandlerJSON(c *fiber.Ctx, err error) error {
	return errorHandler(c, err, true)
}

func errorHandler(c *fiber.Ctx, err error, json bool) error {
	if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreAsSuccess != 0 {
		return c.SendStatus(http.StatusOK)
	}

	status, message := statusAndMessage(err)

	c.Status(status)
	switch {
	case c.Method() == http.MethodHead:
		return nil
	case json:
		return c.JSON(response.Error{HttpStatusCode: status, ErrorCode: strconv.Itoa(status), Message: message})
	default:
		return c.SendString(message)
	}
}

// statusAndMessage reads the status fiber attached to err (a *fiber.Error,
// or a sentinel like fiber.ErrNotFound) and falls back to ErrorStatusCode.
func statusAndMessage(err error) (int, string) {
	var fe *fiber.Error
	if stderrors.As(err, &fe) {
		return fe.Code, fe.Message
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
