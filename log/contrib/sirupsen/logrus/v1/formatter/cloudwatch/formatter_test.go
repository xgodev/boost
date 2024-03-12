package cloudwatch

import (
	//"github.com/stretchr/testify/mock"
	"github.com/ravernkoh/cwlogsfmt"

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

func buildBasicFormatterForTesting() *cwlogsfmt.CloudWatchLogsFormatter {
	return &cwlogsfmt.CloudWatchLogsFormatter{
		PrefixFields:     []string{"RequestId"},
		DisableSorting:   false,
		QuoteEmptyFields: true,
	}
}

func (s *FormatterSuite) TestNewFormatter() {

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
			name: "New Formatter with prefix fields",
			want: func() logrus.Formatter {
				fmt := buildBasicFormatterForTesting()
				WithPrefixFields([]string{"EventParentId"})(fmt)
				return fmt
			},
			opts: []Option{
				WithPrefixFields([]string{"EventParentId"}),
			},
		},
		{
			name: "New Formatter with disabled sorting",
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
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := New(t.opts...)
			want := t.want()
			s.Assert().True(reflect.DeepEqual(got, want), "got  %v\nwant %v", got, want)
		})
	}
}
