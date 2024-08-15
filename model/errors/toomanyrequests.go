package errors

// tooManyRequests represents an error when a request cannot be completed because of
// sent too many requests in a given amount of  time
type tooManyRequests struct {
	Err
}

// TooManyRequestsf returns an error which satistifes IsTooManyRequests()
func TooManyRequestsf(format string, args ...interface{}) error {
	return &tooManyRequests{wrap(nil, format, "", args...)}
}

// NewTooManyRequests returns an error which wraps err that satisfies
// IsTooManyRequests().
func NewTooManyRequests(err error, msg string) error {
	return &tooManyRequests{wrap(err, msg, "")}
}

// IsTooManyRequests reports whether err was created with TooManyRequestsf() or
// NewTooManyRequests().
func IsTooManyRequests(err error) bool {
	err = Cause(err)
	_, ok := err.(*tooManyRequests)
	return ok
}
