package config

import (
	"github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EntriesSuite struct {
	suite.Suite
}

func TestEntriesSuite(t *testing.T) {
	suite.Run(t, new(EntriesSuite))
}

func (s *WrapperSuite) TestAdd() {

	tt := []struct {
		name string
		init func()
		got  func() interface{}
		want interface{}
	}{
		{
			name: "Add without options",
			init: func() { Add("app.application.string", "app_test", "description of string") },
			got:  func() interface{} { return String("app.application.string") },
			want: "app_test",
		},
		{
			name: "Add with options",
			init: func() { Add("app.application.string", "app_test", "description of string", WithHide()) },
			got:  func() interface{} { return String("app.application.string") },
			want: "app_test",
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			t.init()
			Load()
			got := t.got()
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *EntriesSuite) TestEntries() {

	tt := []struct {
		name string
		got  func() Config
		want func() Config
	}{
		{
			name: "Add without options",
			got: func() Config {
				Add("app.application.string", "app_test", "description of string")
				config := Entries()[0]
				return config
			},
			want: func() Config {
				return Config{
					Key:         "app.application.string",
					Value:       "app_test",
					Description: "description of string",
					Options:     &Options{},
				}
			},
		},
		{
			name: "Add with options",
			got: func() Config {
				Add("app.application.string", "app_test", "description of string", WithHide())
				config := Entries()[1]
				return config
			},
			want: func() Config {
				return Config{
					Key:         "app.application.string",
					Value:       "app_test",
					Description: "description of string",
					Options: &Options{
						Hide: true,
					},
				}
			},
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			koanf.New()
			Load()
			got := t.got()
			want := t.want()
			s.Assert().True(reflect.DeepEqual(got, want), "got  %v\nwant %v", got, want)
		})
	}
}
