package errors

import "testing"

func TestBuiltinKind(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want Kind
	}{
		{"NotFound", NotFoundf("x"), KindNotFound},
		{"BadRequest", BadRequestf("x"), KindBadRequest},
		{"NotValid", NotValidf("x"), KindNotValid},
		{"Conflict", Conflictf("x"), KindConflict},
		{"AlreadyExists", AlreadyExistsf("x"), KindAlreadyExists},
		{"Forbidden", Forbiddenf("x"), KindForbidden},
		{"Unauthorized", Unauthorizedf("x"), KindUnauthorized},
		{"ServiceUnavailable", ServiceUnavailablef("x"), KindServiceUnavailable},
		{"NotImplemented", NotImplementedf("x"), KindNotImplemented},
		{"NotProvisioned", NotProvisionedf("x"), KindNotProvisioned},
		{"NotSupported", NotSupportedf("x"), KindNotSupported},
		{"NotAssigned", NotAssignedf("x"), KindNotAssigned},
		{"MethodNotAllowed", MethodNotAllowedf("x"), KindMethodNotAllowed},
		{"TooManyRequests", TooManyRequestsf("x"), KindTooManyRequests},
		{"Timeout", Timeoutf("x"), KindTimeout},
		{"Internal", Internalf("x"), KindInternal},
		{"PlainNew", New("x"), KindInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := builtinKind(tt.err); got != tt.want {
				t.Fatalf("builtinKind(%s) = %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}
