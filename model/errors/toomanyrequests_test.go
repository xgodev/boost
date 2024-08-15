package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TooManyRequestsSuite struct {
	suite.Suite
}

func TestTooManyRequestsSuite(t *testing.T) {
	suite.Run(t, new(TooManyRequestsSuite))
}

func (s *TooManyRequestsSuite) TestTooManyRequestsf() {

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
			name: "TooManyRequests f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: TooManyRequestsf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := TooManyRequestsf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *TooManyRequestsSuite) TestNewTooManyRequests() {

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
			name: "New TooManyRequests",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewTooManyRequests(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewTooManyRequests(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *TooManyRequestsSuite) TestIsTooManyRequests() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is TooManyRequests",
			got:  func() bool { return IsTooManyRequests(TooManyRequestsf("test", "123")) },
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
