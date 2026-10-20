package errors

// notSupported represents an error when something is not supported.
type notSupported struct {
	Err
}

// NotSupportedf returns an error which satisfies IsNotSupported().
// Its message is the formatted text followed by " not supported", so the format
// names the thing: NotSupportedf("sku %s", sku) says "sku cap not supported".
func NotSupportedf(format string, args ...interface{}) error {
	return &notSupported{wrap(nil, format, " not supported", args...)}
}

// NewNotSupported returns an error which wraps err and satisfies
// IsNotSupported().
func NewNotSupported(err error, msg string) error {
	return &notSupported{wrap(err, msg, "")}
}

// IsNotSupported reports whether the error was created with
// NotSupportedf() or NewNotSupported().
func IsNotSupported(err error) bool {
	err = Cause(err)
	_, ok := err.(*notSupported)
	return ok
}
