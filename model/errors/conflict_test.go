package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConflictSuite struct {
	suite.Suite
}

func TestConflictSuite(t *testing.T) {
	suite.Run(t, new(ConflictSuite))
}

func (s *ConflictSuite) TestConflictf() {

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
			name: "Conflict f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: Conflictf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := Conflictf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *ConflictSuite) TestNewConflict() {

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
			name: "New Conflict",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewConflict(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewConflict(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *ConflictSuite) TestIsConflict() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Conflict",
			got:  func() bool { return IsConflict(Conflictf("test", "123")) },
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
