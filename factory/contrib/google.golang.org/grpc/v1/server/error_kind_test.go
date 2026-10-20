package server

import (
	"testing"

	"github.com/xgodev/boost/model/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestError_ConflictIsNotADuplicate: a conflict (the state does not allow
// the call, HTTP 409) answers Aborted; only a duplicate answers
// AlreadyExists.
func TestError_ConflictIsNotADuplicate(t *testing.T) {
	for _, c := range []struct {
		name string
		err  error
		want codes.Code
	}{
		{"conflict", errors.Conflictf("sold out"), codes.Aborted},
		{"already exists", errors.AlreadyExistsf("duplicate"), codes.AlreadyExists},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := status.Code(Error(c.err)); got != c.want {
				t.Fatalf("code = %v, want %v", got, c.want)
			}
		})
	}
}
