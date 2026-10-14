package errors_test

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/xgodev/boost/model/errors"
)

// A boost typed error must be transparent to stdlib errors.Is/As so that a
// wrapped sentinel (e.g. context.Canceled) is still discoverable after being
// classified — boost's own Cause() is single-level and does not provide this.
func TestUnwrap_StdlibErrorsIsTraversesBoostError(t *testing.T) {
	err := errors.NewServiceUnavailable(context.Canceled, "vtex down")

	if !errors.IsServiceUnavailable(err) {
		t.Fatalf("expected IsServiceUnavailable to stay true")
	}
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("stderrors.Is must traverse the boost error to context.Canceled")
	}
}

// errors.As must reach a wrapped concrete error type through a boost wrapper.
func TestUnwrap_StdlibErrorsAsReachesWrapped(t *testing.T) {
	sentinel := &customErr{"boom"}
	err := errors.NewInternal(sentinel, "internal failure")

	var target *customErr
	if !stderrors.As(err, &target) {
		t.Fatalf("stderrors.As must reach *customErr through the boost wrapper")
	}
	if target.s != "boom" {
		t.Fatalf("unexpected target: %q", target.s)
	}
}

type customErr struct{ s string }

func (e *customErr) Error() string { return e.s }
