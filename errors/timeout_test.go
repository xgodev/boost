package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TimeoutSuite struct {
	suite.Suite
}

func TestTimeoutSuite(t *testing.T) {
	suite.Run(t, new(TimeoutSuite))
}

func (s *TimeoutSuite) TestTimeoutf() {

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
			name: "Timeout f",
			args: args{
				msg:   msg,
				param: param,
			},
			want: Timeoutf(msg, param),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := Timeoutf(t.args.msg, t.args.param)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *TimeoutSuite) TestNewTimeout() {

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
			name: "New Timeout",
			args: args{
				err: err,
				msg: msg,
			},
			want: NewTimeout(err, msg),
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewTimeout(t.args.err, t.args.msg)
			s.Assert().True(reflect.DeepEqual(got.Error(), t.want.Error()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *TimeoutSuite) TestIsTimeout() {

	tt := []struct {
		name string
		got  func() bool
		want bool
	}{
		{
			name: "Is Timeout",
			got:  func() bool { return IsTimeout(Timeoutf("test", "123")) },
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
