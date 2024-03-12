package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceUnavailableSuite struct {
	suite.Suite
}

func TestServiceUnavailableSuite(t *testing.T) {
	suite.Run(t, new(ServiceUnavailableSuite))
}

func (s *ServiceUnavailableSuite) TestServiceUnavailablef() {

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
			name: "Service Unavailable f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: ServiceUnavailablef(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := ServiceUnavailablef(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *ServiceUnavailableSuite) TestNewServiceUnavailable() {

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
			name: "New Service Unavailable",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewServiceUnavailable(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewServiceUnavailable(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *ServiceUnavailableSuite) TestIsServiceUnavailable() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Service Unavailable",
			got:  func() bool { return IsServiceUnavailable(ServiceUnavailablef("test", "123")) },
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
