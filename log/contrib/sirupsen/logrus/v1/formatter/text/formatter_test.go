package text

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

type FormatterSuite struct {
	suite.Suite
}

func TestFormatterSuite(t *testing.T) {
	suite.Run(t, new(FormatterSuite))
}

func buildBasicFormatterForTesting() *logrus.TextFormatter {
	return &logrus.TextFormatter{
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
}
func (s *FormatterSuite) TestNew() {
	//callerPrettyfier := func(rtf *runtime.Frame) (function string, file string) { return "", "" }
	tt := []struct {
		name string
		want func() logrus.Formatter
		opts []Option
	}{
		{
			name: "New Formatter with default options",
			want: func() logrus.Formatter {
				return buildBasicFormatterForTesting()
			},
			opts: []Option{},
		},
		{
			name: "New Formatter with caller prettyfier",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithCallerPrettyfier(nil)(fmt)
				//Note deepEqual and func https://github.com/golang/go/issues/8554
				return fmt
			},
			opts: []Option{
				WithCallerPrettyfier(nil),
			},
		},
		{
			name: "New Formatter with disable colors",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithDisableColors(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithDisableColors(true),
			},
		},
		{
			name: "New Formatter with disable level truncation",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithDisableLevelTruncation(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithDisableLevelTruncation(true),
			},
		},
		{
			name: "New Formatter with disable quote",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithDisableQuote(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithDisableQuote(true),
			},
		},
		{
			name: "New Formatter with disable sorting",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithDisableSorting(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithDisableSorting(true),
			},
		},
		{
			name: "New Formatter with disable timestamp",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithDisableTimestamp(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithDisableTimestamp(true),
			},
		},
		{
			name: "New Formatter with environment override colors",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithEnvironmentOverrideColors(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithEnvironmentOverrideColors(true),
			},
		},
		{
			name: "New Formatter with field map",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithFieldMap(logrus.FieldMap{logrus.FieldKeyLevel: "TRACE"})(fmt)
				return fmt
			},
			opts: []Option{
				WithFieldMap(logrus.FieldMap{logrus.FieldKeyLevel: "TRACE"}),
			},
		},
		{
			name: "New Formatter with force colors",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithForceColors(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithForceColors(true),
			},
		},
		{
			name: "New Formatter with force Quote",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithForceQuote(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithForceQuote(true),
			},
		},
		{
			name: "New Formatter with full timestamp",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithFullTimestamp(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithFullTimestamp(true),
			},
		},
		{
			name: "New Formatter with pad level text",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithPadLevelText(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithPadLevelText(true),
			},
		},
		{
			name: "New Formatter with quote empty fields",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithQuoteEmptyFields(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithQuoteEmptyFields(true),
			},
		},
		{
			name: "New Formatter with sorting func",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithSortingFunc(nil)(fmt)
				//Note deepEqual and func https://github.com/golang/go/issues/8554
				return fmt
			},
			opts: []Option{
				WithSortingFunc(nil),
			},
		},
		{
			name: "New Formatter with timestamp format",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithTimestampFormat("2006")(fmt)
				return fmt
			},
			opts: []Option{
				WithTimestampFormat("2006"),
			},
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := New(t.opts...)
			want := t.want()
			s.Assert().True(reflect.DeepEqual(got, want), "\ngot  %v\nwant %v", got, want)
		})
	}
}
