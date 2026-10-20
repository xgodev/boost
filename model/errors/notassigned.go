package errors

// notAssigned represents an error when something is not yet assigned to
// something else.
type notAssigned struct {
	Err
}

// NotAssignedf returns an error which satisfies IsNotAssigned().
// Its message is the formatted text followed by " not assigned", so the format
// names the thing: NotAssignedf("sku %s", sku) says "sku cap not assigned".
func NotAssignedf(format string, args ...interface{}) error {
	return &notAssigned{wrap(nil, format, " not assigned", args...)}
}

// NewNotAssigned returns an error which wraps err that satisfies
// IsNotAssigned().
func NewNotAssigned(err error, msg string) error {
	return &notAssigned{wrap(err, msg, "")}
}

// IsNotAssigned reports whether err was created with NotAssignedf() or
// NewNotAssigned().
func IsNotAssigned(err error) bool {
	err = Cause(err)
	_, ok := err.(*notAssigned)
	return ok
}
