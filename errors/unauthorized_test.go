package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UnauthorizedSuite struct {
	suite.Suite
}

func TestUnauthorizedSuite(t *testing.T) {
	suite.Run(t, new(UnauthorizedSuite))
}

func (s *UnauthorizedSuite) TestUnauthorizedf() {

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
			name: "Unauthorized f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: Unauthorizedf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := Unauthorizedf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *UnauthorizedSuite) TestNewUnauthorized() {

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
			name: "New Unauthorized",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewUnauthorized(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewUnauthorized(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *UnauthorizedSuite) TestIsUnauthorized() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Unauthorized",
			got:  func() bool { return IsUnauthorized(Unauthorizedf("test", "123")) },
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
