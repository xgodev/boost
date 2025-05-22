package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotFoundSuite struct {
	suite.Suite
}

func TestNotFoundSuite(t *testing.T) {
	suite.Run(t, new(NotFoundSuite))
}

func (s *NotFoundSuite) TestNotFoundf() {

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
			name: "Not Found f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: NotFoundf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NotFoundf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotFoundSuite) TestNewNotFound() {

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
			name: "New Not Found",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewNotFound(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewNotFound(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotFoundSuite) TestIsNotFound() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Not Found",
			got:  func() bool { return IsNotFound(NotFoundf("test", "123")) },
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
