package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotSupportedSuite struct {
	suite.Suite
}

func TestNotSupportedSuite(t *testing.T) {
	suite.Run(t, new(NotSupportedSuite))
}

func (s *NotSupportedSuite) TestNotSupportedf() {

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
			name: "Not Supported f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: NotSupportedf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NotSupportedf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotSupportedSuite) TestNewNotSupported() {

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
			name: "New Not Supported",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewNotSupported(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewNotSupported(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotSupportedSuite) TestIsNotSupported() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Not Supported",
			got:  func() bool { return IsNotSupported(NotSupportedf("test", "123")) },
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
