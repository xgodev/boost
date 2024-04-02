package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NotAssignedSuite struct {
	suite.Suite
}

func TestNotAssignedSuite(t *testing.T) {
	suite.Run(t, new(NotAssignedSuite))
}

func (s *NotAssignedSuite) TestNotAssignedf() {

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
			name: "Not Assigned f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: NotAssignedf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NotAssignedf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotAssignedSuite) TestNewNotAssigned() {

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
			name: "New Not Assigned",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewNotAssigned(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewNotAssigned(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *NotAssignedSuite) TestIsNotAssigned() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Not Assigned",
			got:  func() bool { return IsNotAssigned(NotAssignedf("test", "123")) },
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
