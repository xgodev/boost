package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotProvisionedSuite struct {
	suite.Suite
}

func TestNotProvisionedSuite(t *testing.T) {
	suite.Run(t, new(NotProvisionedSuite))
}

func (s *NotProvisionedSuite) TestNotProvisionedf() {

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
			name: "Not Provisioned f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: NotProvisionedf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NotProvisionedf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotProvisionedSuite) TestNewNotProvisioned() {

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
			name: "New Not Provisioned",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewNotProvisioned(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewNotProvisioned(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotProvisionedSuite) TestIsNotProvisioned() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Not Provisioned",
			got:  func() bool { return IsNotProvisioned(NotProvisionedf("test", "123")) },
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
