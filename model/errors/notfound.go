package errors

// notFound represents an error when something has not been found.
type notFound struct {
	Err
}

// NotFoundf returns an error which satisfies IsNotFound().
// Its message is the formatted text followed by " not found", so the format
// names the thing: NotFoundf("sku %s", sku) says "sku cap not found".
func NotFoundf(format string, args ...interface{}) error {
	return &notFound{wrap(nil, format, " not found", args...)}
}

// NewNotFound returns an error which wraps err that satisfies
// IsNotFound().
func NewNotFound(err error, msg string) error {
	return &notFound{wrap(err, msg, "")}
}

// IsNotFound reports whether err was created with NotFoundf() or
// NewNotFound().
func IsNotFound(err error) bool {
	err = Cause(err)
	_, ok := err.(*notFound)
	return ok
}
