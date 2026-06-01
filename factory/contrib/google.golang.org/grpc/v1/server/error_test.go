package server

import (
	"testing"

	"github.com/xgodev/boost/model/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type customConflict struct{ msg string }

func (c *customConflict) Error() string { return c.msg }

func TestError_RegisteredCustom(t *testing.T) {
	errors.RegisterMatch(func(err error) bool {
		_, ok := err.(*customConflict)
		return ok
	}, errors.KindConflict)

	if got := status.Code(Error(&customConflict{"dup"})); got != codes.AlreadyExists {
		t.Fatalf("code = %v, want AlreadyExists", got)
	}
}

func TestError_IgnoreAsSuccess(t *testing.T) {
	sentinel := errors.New("ignorable-grpc")
	errors.Ignore(sentinel, errors.IgnoreAsSuccess)

	if err := Error(sentinel); err != nil {
		t.Fatalf("Error = %v, want nil", err)
	}
}

func TestError_BuiltinNotFound(t *testing.T) {
	if status.Code(Error(errors.NotFoundf("x"))) != codes.NotFound {
		t.Fatal("builtin NotFound not mapped to codes.NotFound")
	}
}

func TestError_DefaultInternal(t *testing.T) {
	if status.Code(Error(errors.New("plain"))) != codes.Internal {
		t.Fatal("plain error not mapped to codes.Internal")
	}
}
