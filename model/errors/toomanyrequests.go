package errors

// toManyRequests represents an error when a request cannot be completed because of
// sent too many requests in a given amount of  time
type toManyRequests struct {
	Err
}

// ToManyRequestsf returns an error which satistifes IsToManyRequests()
func ToManyRequestsf(format string, args ...interface{}) error {
	return &toManyRequests{wrap(nil, format, "", args...)}
}

// NewToManyRequests returns an error which wraps err that satisfies
// IsToManyRequests().
func NewToManyRequests(err error, msg string) error {
	return &toManyRequests{wrap(err, msg, "")}
}

// IsToManyRequests reports whether err was created with ToManyRequestsf() or
// NewInternal().
func IsToManyRequests(err error) bool {
	err = Cause(err)
	_, ok := err.(*toManyRequests)
	return ok
}
