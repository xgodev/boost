package json

import "github.com/sirupsen/logrus"

// Option represents a JSON formatter option.
type Option func(formatter *logrus.JSONFormatter)

// New returns a new logrus formatter for JSON.
func New(options ...Option) logrus.Formatter {
	fmt := &logrus.JSONFormatter{
		TimestampFormat:   "2006/01/02 15:04:05.000",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		PrettyPrint:       false,
	}

	for _, option := range options {
		option(fmt)
	}

	return fmt
}

// WithTimestampFormat sets formatter's timestamp format to value.
// The format to use is the same than for time.Format or time.Parse from the standard library.
func WithTimestampFormat(value string) Option {
	return func(formatter *logrus.JSONFormatter) {
		formatter.TimestampFormat = value
	}
}

// WithDisableTimestamp sets formatter's disable timestamp to value.
func WithDisableTimestamp(value bool) Option {
	return func(formatter *logrus.JSONFormatter) {
		formatter.DisableTimestamp = value
	}
}

// WithDisableHTMLEscape sets formatter's disable HTML escape to value.
func WithDisableHTMLEscape(value bool) Option {
	return func(formatter *logrus.JSONFormatter) {
		formatter.DisableHTMLEscape = value
	}
}

// WithDataKey sets formatter's data key to value.
func WithDataKey(value string) Option {
	return func(formatter *logrus.JSONFormatter) {
		formatter.DataKey = value
	}
}

// WithFieldMap sets formatter's field map to value.
func WithFieldMap(value logrus.FieldMap) Option {
	return func(formatter *logrus.JSONFormatter) {
		formatter.FieldMap = value
	}
}

// WithPrettyPrint sets formatter's pretty print to value.
func WithPrettyPrint(value bool) Option {
	return func(formatter *logrus.JSONFormatter) {
		formatter.PrettyPrint = value
	}
}
