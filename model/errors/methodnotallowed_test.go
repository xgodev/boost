package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MethodNotAllowedSuite struct {
	suite.Suite
}

func TestMethodNotAllowedSuite(t *testing.T) {
	suite.Run(t, new(MethodNotAllowedSuite))
}

func (s *MethodNotAllowedSuite) TestMethodNotAllowedf() {

	msg := "test"
	param := "123"

	type args struct {
		msg   string
		param string
	}

	tt := []struct {
		name string
		args args
		want error
	}{
		{
			name: "Method Not Allowed f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: MethodNotAllowedf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := MethodNotAllowedf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *MethodNotAllowedSuite) TestNewMethodNotAllowed() {

	err := errors.New("test")
	msg := "test"

	type args struct {
		err error
		msg string
	}

	tt := []struct {
		name string
		args args
		want error
	}{
		{
			name: "New Method Not Allowed",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewMethodNotAllowed(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewMethodNotAllowed(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *MethodNotAllowedSuite) TestIsMethodNotAllowed() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Method Not Allowed",
			got:  func() bool { return IsMethodNotAllowed(MethodNotAllowedf("test", "123")) },
			want: true,
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := t.got()
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}
