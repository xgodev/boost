package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type InternalSuite struct {
	suite.Suite
}

func TestInternalSuite(t *testing.T) {
	suite.Run(t, new(InternalSuite))
}

func (s *InternalSuite) TestInternalf() {

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
			name: "Internal f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: Internalf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := Internalf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *InternalSuite) TestNewInternal() {

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
			name: "New Internal",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewInternal(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewInternal(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *InternalSuite) TestIsInternal() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Internal",
			got:  func() bool { return IsInternal(Internalf("test", "123")) },
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
