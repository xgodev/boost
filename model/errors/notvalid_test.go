package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotValidSuite struct {
	suite.Suite
}

func TestNotValidSuite(t *testing.T) {
	suite.Run(t, new(NotValidSuite))
}

func (s *NotValidSuite) TestNotValidf() {

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
			name: "Not Valid f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: NotValidf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NotValidf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotValidSuite) TestNewNotValid() {

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
			name: "New Not Valid",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewNotValid(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewNotValid(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotValidSuite) TestIsNotValid() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Not Valid",
			got:  func() bool { return IsNotValid(NotValidf("test", "123")) },
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
