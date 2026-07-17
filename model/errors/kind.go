package errors

// Kind is a transport-agnostic semantic classification of an error.
type Kind int

const (
	KindInternal Kind = iota // default / fallback
	KindNotFound
	KindBadRequest
	KindNotValid
	KindConflict
	KindAlreadyExists
	KindForbidden
	KindUnauthorized
	KindServiceUnavailable
	KindNotImplemented
	KindNotProvisioned
	KindNotSupported
	KindNotAssigned
	KindMethodNotAllowed
	KindTooManyRequests
	KindTimeout
)

// builtinKind classifies err using the typed error catalog (Is* checks).
// More specific types are checked first; IsNotValid is checked before
// IsBadRequest because the two may overlap. Unmatched errors fall back to
// KindInternal.
func builtinKind(err error) Kind {
	switch {
	case IsNotFound(err):
		return KindNotFound
	case IsMethodNotAllowed(err):
		return KindMethodNotAllowed
	case IsNotValid(err):
		return KindNotValid
	case IsBadRequest(err):
		return KindBadRequest
	case IsServiceUnavailable(err):
		return KindServiceUnavailable
	case IsConflict(err):
		return KindConflict
	case IsAlreadyExists(err):
		return KindAlreadyExists
	case IsNotImplemented(err):
		return KindNotImplemented
	case IsNotProvisioned(err):
		return KindNotProvisioned
	case IsUnauthorized(err):
		return KindUnauthorized
	case IsForbidden(err):
		return KindForbidden
	case IsNotSupported(err):
		return KindNotSupported
	case IsNotAssigned(err):
		return KindNotAssigned
	case IsTooManyRequests(err):
		return KindTooManyRequests
	case IsTimeout(err):
		return KindTimeout
	default:
		return KindInternal
	}
}
