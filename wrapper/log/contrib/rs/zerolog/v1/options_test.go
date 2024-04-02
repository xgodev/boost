package zerolog

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
			name:   "Options with console enabled",
			want:   true,
			got:    func(o *Options) interface{} { return o.Console.Enabled },
			method: WithConsoleEnabled(true),
		},
		{
			name:   "Options with level",
			want:   "TRACE",
			got:    func(o *Options) interface{} { return o.Level },
			method: WithLevel("TRACE"),
		},
		{
			name:   "Options with file compress",
			want:   true,
			got:    func(o *Options) interface{} { return o.File.Compress },
			method: WithFileCompress(true),
		},
		{
			name:   "Options with file enabled",
			want:   true,
			got:    func(o *Options) interface{} { return o.File.Enabled },
			method: WithFileEnabled(true),
		},
		{
			name:   "Options with file max age",
			want:   7,
			got:    func(o *Options) interface{} { return o.File.MaxAge },
			method: WithFileMaxAge(7),
		},
		{
			name:   "Options with file max size",
			want:   50,
			got:    func(o *Options) interface{} { return o.File.MaxSize },
			method: WithFileMaxSize(50),
		},
		{
			name:   "Options with file name",
			want:   "app.log",
			got:    func(o *Options) interface{} { return o.File.Name },
			method: WithFileName("app.log"),
		},
		{
			name:   "Options with file path",
			want:   "/temporary",
			got:    func(o *Options) interface{} { return o.File.Path },
			method: WithFilePath("/temporary"),
		},
		{
			name:   "Options with formatter",
			want:   "JSON",
			got:    func(o *Options) interface{} { return o.Formatter },
			method: WithFormatter("JSON"),
		},
		{
			name:   "Options with error field name",
			want:   "error",
			got:    func(o *Options) interface{} { return o.ErrorFieldName },
			method: WithErrorFieldName("error"),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			opts := defaultOptions()
			t.method(opts)
			got := t.got(opts)
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}
