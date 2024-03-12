package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AlreadyExistsSuite struct {
	suite.Suite
}

func TestAlreadyExistsSuite(t *testing.T) {
	suite.Run(t, new(AlreadyExistsSuite))
}

func (s *AlreadyExistsSuite) TestAlreadyExistsf() {

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
			name: "Already Exists f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: AlreadyExistsf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := AlreadyExistsf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *AlreadyExistsSuite) TestNewAlreadyExists() {

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
			name: "New Already Exists",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewAlreadyExists(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewAlreadyExists(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *AlreadyExistsSuite) TestIsAlreadyExists() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Already Exists",
			got:  func() bool { return IsAlreadyExists(AlreadyExistsf("test", "123")) },
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
