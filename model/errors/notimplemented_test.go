package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotImplementedSuite struct {
	suite.Suite
}

func TestNotImplementedSuite(t *testing.T) {
	suite.Run(t, new(NotImplementedSuite))
}

func (s *NotImplementedSuite) TestNotImplementedf() {

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
			name: "Not Implemented f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: NotImplementedf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NotImplementedf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotImplementedSuite) TestNewNotImplemented() {

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
			name: "New Not Implemented",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewNotImplemented(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewNotImplemented(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotImplementedSuite) TestIsNotImplemented() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Not Implemented",
			got:  func() bool { return IsNotImplemented(NotImplementedf("test", "123")) },
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
