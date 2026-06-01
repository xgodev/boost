package http

import (
	"net/http"
	"testing"

	"github.com/xgodev/boost/model/errors"
)

type customForbidden struct{ msg string }

func (c *customForbidden) Error() string { return c.msg }

func TestErrorStatusCode_RegisteredCustom(t *testing.T) {
	errors.RegisterMatch(func(err error) bool {
		_, ok := err.(*customForbidden)
		return ok
	}, errors.KindForbidden)

	if got := ErrorStatusCode(&customForbidden{"no"}); got != http.StatusForbidden {
		t.Fatalf("ErrorStatusCode = %d, want 403", got)
	}
}

func TestErrorStatusCode_BuiltinNotFound(t *testing.T) {
	if ErrorStatusCode(errors.NotFoundf("x")) != http.StatusNotFound {
		t.Fatal("builtin NotFound not mapped to 404")
	}
}
