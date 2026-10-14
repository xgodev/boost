package response

import (
	"net/http"

	"github.com/xgodev/boost/model/errors"
)

// HTTPStatusFor maps a boost error Kind to its HTTP status code.
func HTTPStatusFor(kind errors.Kind) int {
	switch kind {
	case errors.KindNotFound:
		return http.StatusNotFound
	case errors.KindMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case errors.KindNotValid, errors.KindBadRequest:
		return http.StatusBadRequest
	case errors.KindServiceUnavailable:
		return http.StatusServiceUnavailable
	case errors.KindConflict, errors.KindAlreadyExists:
		return http.StatusConflict
	case errors.KindNotImplemented, errors.KindNotProvisioned:
		return http.StatusNotImplemented
	case errors.KindUnauthorized:
		return http.StatusUnauthorized
	case errors.KindForbidden:
		return http.StatusForbidden
	case errors.KindNotSupported, errors.KindNotAssigned:
		return http.StatusUnprocessableEntity
	case errors.KindTooManyRequests:
		return http.StatusTooManyRequests
	case errors.KindTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}
