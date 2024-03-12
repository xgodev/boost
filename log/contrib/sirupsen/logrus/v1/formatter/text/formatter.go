package text

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

// Option represents a text formatter options.
type Option func(formatter *logrus.TextFormatter)

// New returns a new logrus formatter for text.
func New(options ...Option) logrus.Formatter {

	fmt := &logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "2006/01/02 15:04:05.000",
		DisableSorting:            false,
		DisableLevelTruncation:    true,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
	}

	for _, option := range options {
		option(fmt)
	}

	return fmt
}

// WithForceColors sets formatter's force colors to value.
func WithForceColors(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.ForceColors = value
	}
}

// WithDisableColors sets formatter's disable colors to value.
func WithDisableColors(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.DisableColors = value
	}
}

// WithForceQuote sets formatter's force quote to value.
func WithForceQuote(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.ForceQuote = value
	}
}

// WithDisableQuote sets formatter's disable quote to value.
func WithDisableQuote(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.DisableQuote = value
	}
}

// WithEnviromentOverrideColors sets formatter's override colors to value.
// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
func WithEnvironmentOverrideColors(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.EnvironmentOverrideColors = value
	}
}

// WithDisableTimestamp sets formatter's disable timestamp to value.
// Disable timestamps is useful when output is redirected to logging system that already adds timestamps.
func WithDisableTimestamp(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.DisableTimestamp = value
	}
}

// WithFullTimestamp sets formatter's full timestamp to value.
func WithFullTimestamp(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.FullTimestamp = value
	}
}

// WithTimestampFormat sets formatter's timestamp format to value.
// The format to use is the same than for time.Format or time.Parse from the standard library.
func WithTimestampFormat(value string) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.TimestampFormat = value
	}
}

// WithDisableSorting sets formatter's disable sorting to value.
// The fields are sorted by default for a consistent output. For applications that log extremely frequently and don't use the JSON formatter this may not be desired.
func WithDisableSorting(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.DisableSorting = value
	}
}

// WithDisableLevelTruncation sets formatter's disable level truncation to value.
// Disables the truncation of the level text to 4 characters.
func WithDisableLevelTruncation(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.DisableLevelTruncation = value
	}
}

// WithPadLevelText sets formatter's pad level text to value.
// Adds padding the level text so that all the levels output at the same length PadLevelText is a superset of the DisableLevelTruncation option.
func WithPadLevelText(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.PadLevelText = value
	}
}

// WithQuoteEmptyFields sets formatter's quote empty fields to value.
// This will wrap empty fields in quotes if true.
func WithQuoteEmptyFields(value bool) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.QuoteEmptyFields = value
	}
}

// WithFieldMap sets formatter's field map to value.
func WithFieldMap(value logrus.FieldMap) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.FieldMap = value
	}
}

// WithCallerPrettyfier sets formatter's caller prettyfier to value.
// Can be set by the user to modify the content of the function and file keys in the data when ReportCaller is activated. If any of the returned value is the empty string the corresponding key will be removed from fields.
func WithCallerPrettyfier(value func(*runtime.Frame) (function string, file string)) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.CallerPrettyfier = value
	}
}

// WithSortingFunc sets formatter's sorting func to value.
// The keys sorting function, when uninitialized it uses sort.Strings.
func WithSortingFunc(value func([]string)) Option {
	return func(formatter *logrus.TextFormatter) {
		formatter.SortingFunc = value
	}
}
