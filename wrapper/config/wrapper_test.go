package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type WrapperSuite struct {
	suite.Suite
}

func TestWrapperSuite(t *testing.T) {
	suite.Run(t, new(WrapperSuite))
}

type AppConfig struct {
	Application struct {
		Name string
	}
}

func (s *WrapperSuite) TestWrapperWrappedMethods() {

	tt := []struct {
		name string
		init func()
		got  func() interface{}
		want interface{}
	}{
		{
			name: "UnmarshalWithPath",
			init: func() {},
			got:  func() interface{} { return UnmarshalWithPath("app", &AppConfig{}) },
			want: nil,
		},
		{
			name: "Unmarshal",
			init: func() {},
			got:  func() interface{} { return Unmarshal(&AppConfig{}) },
			want: nil,
		},
		{
			name: "Exists",
			init: func() { Add("app.application.exists", true, "description of exists") },
			got:  func() interface{} { return Exists("app.application.exists") },
			want: true,
		},
		{
			name: "Int64",
			init: func() { Add("app.application.int64", int64(10), "description of int64") },
			got:  func() interface{} { return Int64("app.application.int64") },
			want: int64(10),
		},
		// {
		// 	name: "Int64s",
		// 	init: func() { Add("app.application.durations", []int64{8, 9, 10}, "description of durations") },
		// 	got:  func() interface{} { return Int64s("app.application.durations") },
		// 	want: []int64{8, 9, 10},
		// },
		{
			name: "Int",
			init: func() { Add("app.application.int", 10, "description of int") },
			got:  func() interface{} { return Int("app.application.int") },
			want: 10,
		},
		{
			name: "Ints",
			init: func() { Add("app.application.ints", []int{8, 9, 10}, "description of ints") },
			got:  func() interface{} { return Ints("app.application.ints") },
			want: []int{8, 9, 10},
		},
		{
			name: "Float64",
			init: func() { Add("app.application.float64", float64(2.55), "description of float64") },
			got:  func() interface{} { return Float64("app.application.float64") },
			want: float64(2.55),
		},
		{
			name: "Duration",
			init: func() { Add("app.application.duration", time.Second, "description of duration") },
			got:  func() interface{} { return Duration("app.application.duration") },
			want: time.Second,
		},
		{
			name: "String",
			init: func() { Add("app.application.string", "app_test", "description of string") },
			got:  func() interface{} { return String("app.application.string") },
			want: "app_test",
		},
		{
			name: "Strings",
			init: func() { Add("app.application.strings", []string{"app_test", "app_test2"}, "description of strings") },
			got:  func() interface{} { return Strings("app.application.strings") },
			want: []string{"app_test", "app_test2"},
		},
		{
			name: "Bytes",
			init: func() { Add("app.application.bytes", "test", "description of bytes") },
			got:  func() interface{} { return Bytes("app.application.bytes") },
			want: []byte{116, 101, 115, 116}, //test in base64 array
		},
		{
			name: "Bool",
			init: func() { Add("app.application.bool", true, "description of bool") },
			got:  func() interface{} { return Bool("app.application.bool") },
			want: true,
		},
		// {
		// 	name: "Bools",
		// 	init: func() { Add("app.application.bools", []bool{true, false}, "description of bools") },
		// 	got:  func() interface{} { return Bools("app.application.bools") },
		// 	want: []bool{true, false},
		// },
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
