package cloudwatch

import (
	"github.com/ravernkoh/cwlogsfmt"
	"github.com/sirupsen/logrus"
)

// Option represents a CloudWatch formatter option.
type Option func(formatter *cwlogsfmt.CloudWatchLogsFormatter)

// New returns a new logrus formatter for CloudWatch.
func New(options ...Option) logrus.Formatter {
	fmt := &cwlogsfmt.CloudWatchLogsFormatter{
		PrefixFields:     []string{"RequestId"},
		DisableSorting:   false,
		QuoteEmptyFields: true,
	}

	for _, option := range options {
		option(fmt)
	}

	return fmt
}

// WithPrefixFields sets formatter's prefix fields to value.
func WithPrefixFields(value []string) Option {
	return func(formatter *cwlogsfmt.CloudWatchLogsFormatter) {
		formatter.PrefixFields = value
	}
}

// WithDisableSorting sets formatter's disable sorting to value.
func WithDisableSorting(value bool) Option {
	return func(formatter *cwlogsfmt.CloudWatchLogsFormatter) {
		formatter.DisableSorting = value
	}
}

// WithQuoteEmptyFields sets formatter's quote empty fields to value.
func WithQuoteEmptyFields(value bool) Option {
	return func(formatter *cwlogsfmt.CloudWatchLogsFormatter) {
		formatter.QuoteEmptyFields = value
	}
}
