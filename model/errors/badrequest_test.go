package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BadRequestSuite struct {
	suite.Suite
}

func TestBadRequestSuite(t *testing.T) {
	suite.Run(t, new(BadRequestSuite))
}

func (s *BadRequestSuite) TestBadRequestf() {

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
			name: "Bad Request f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: BadRequestf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := BadRequestf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *BadRequestSuite) TestNewBadRequest() {

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
			name: "New Bad Request",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewBadRequest(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewBadRequest(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *BadRequestSuite) TestIsBadRequest() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Bad Request",
			got:  func() bool { return IsBadRequest(BadRequestf("test", "123")) },
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
