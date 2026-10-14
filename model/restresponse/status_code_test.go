package response

import (
	"net/http"
	"testing"

	"github.com/xgodev/boost/model/errors"
)

func TestHTTPStatusFor(t *testing.T) {
	tests := []struct {
		name string
		kind errors.Kind
		want int
	}{
		{"NotFound", errors.KindNotFound, http.StatusNotFound},
		{"MethodNotAllowed", errors.KindMethodNotAllowed, http.StatusMethodNotAllowed},
		{"NotValid", errors.KindNotValid, http.StatusBadRequest},
		{"BadRequest", errors.KindBadRequest, http.StatusBadRequest},
		{"ServiceUnavailable", errors.KindServiceUnavailable, http.StatusServiceUnavailable},
		{"Conflict", errors.KindConflict, http.StatusConflict},
		{"AlreadyExists", errors.KindAlreadyExists, http.StatusConflict},
		{"NotImplemented", errors.KindNotImplemented, http.StatusNotImplemented},
		{"NotProvisioned", errors.KindNotProvisioned, http.StatusNotImplemented},
		{"Unauthorized", errors.KindUnauthorized, http.StatusUnauthorized},
		{"Forbidden", errors.KindForbidden, http.StatusForbidden},
		{"NotSupported", errors.KindNotSupported, http.StatusUnprocessableEntity},
		{"NotAssigned", errors.KindNotAssigned, http.StatusUnprocessableEntity},
		{"TooManyRequests", errors.KindTooManyRequests, http.StatusTooManyRequests},
		{"Timeout", errors.KindTimeout, http.StatusRequestTimeout},
		{"Internal", errors.KindInternal, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HTTPStatusFor(tt.kind); got != tt.want {
				t.Errorf("HTTPStatusFor(%v) = %d, want %d", tt.kind, got, tt.want)
			}
		})
	}
}
