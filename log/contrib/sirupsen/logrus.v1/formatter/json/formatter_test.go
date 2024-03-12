package json

import (
	//"github.com/stretchr/testify/mock"

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

func buildBasicFormatterForTesting() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat:   "2006/01/02 15:04:05.000",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		PrettyPrint:       false,
	}
}
func (s *FormatterSuite) TestNew() {

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
			name: "New Formatter with data key",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithDataKey("key")(fmt)
				return fmt
			},
			opts: []Option{
				WithDataKey("key"),
			},
		},
		{
			name: "New Formatter with disable HTML escape",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithDisableHTMLEscape(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithDisableHTMLEscape(true),
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
			name: "New Formatter with pretty print",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithPrettyPrint(true)(fmt)
				return fmt
			},
			opts: []Option{
				WithPrettyPrint(true),
			},
		},
		{
			name: "New Formatter with pretty print",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithTimestampFormat("2006-01-02T15:04:05.000")(fmt)
				return fmt
			},
			opts: []Option{
				WithTimestampFormat("2006-01-02T15:04:05.000"),
			},
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := New(t.opts...)
			want := t.want()
			s.Assert().True(reflect.DeepEqual(got, want), "got  %v\nwant %v", got, want)
		})
	}
}
