package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type OptionsSuite struct {
	suite.Suite
}

func TestOptionsSuite(t *testing.T) {
	suite.Run(t, new(OptionsSuite))
}

func (s *OptionsSuite) TestOptionsWithMethods() {

	tt := []struct {
		name   string
		want   interface{}
		got    func(o *Options) interface{}
		method Option
	}{
		{
			name:   "Options with hide true",
			want:   true,
			got:    func(o *Options) interface{} { return o.Hide },
			method: WithHide(),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			opts := &Options{}
			t.method(opts)
			got := t.got(opts)
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}
