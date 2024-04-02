package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ForbiddenSuite struct {
	suite.Suite
}

func TestForbiddenSuite(t *testing.T) {
	suite.Run(t, new(ForbiddenSuite))
}

func (s *ForbiddenSuite) TestForbiddenf() {

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
			name: "Forbidden f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: Forbiddenf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := Forbiddenf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *ForbiddenSuite) TestNewForbidden() {

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
			name: "New Forbidden",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewForbidden(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewForbidden(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *ForbiddenSuite) TestIsForbidden() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Forbidden",
			got:  func() bool { return IsForbidden(Forbiddenf("test", "123")) },
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
