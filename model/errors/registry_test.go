package errors

import (
	stderrors "errors"
	"testing"
)

func TestClassifyRegisteredSentinel(t *testing.T) {
	defer resetRegistry()

	sentinel := New("sentinel not found")
	Register(sentinel, KindNotFound)

	if got := Classify(sentinel); got != KindNotFound {
		t.Fatalf("Classify(sentinel) = %d, want %d", got, KindNotFound)
	}
}

type xptoError struct{ msg string }

func (e *xptoError) Error() string { return e.msg }

func TestClassifyRegisteredMatch(t *testing.T) {
	defer resetRegistry()

	RegisterMatch(func(err error) bool {
		var target *xptoError
		return stderrors.As(err, &target)
	}, KindNotFound)

	if got := Classify(&xptoError{msg: "boom"}); got != KindNotFound {
		t.Fatalf("Classify(xptoError) = %d, want %d", got, KindNotFound)
	}
}

func TestClassifyFallsBackToBuiltin(t *testing.T) {
	defer resetRegistry()

	if got := Classify(NotFoundf("x")); got != KindNotFound {
		t.Fatalf("Classify(NotFoundf) = %d, want %d", got, KindNotFound)
	}
}

func TestClassifyDefaultInternal(t *testing.T) {
	defer resetRegistry()

	if got := Classify(New("anything")); got != KindInternal {
		t.Fatalf("Classify(New) = %d, want %d", got, KindInternal)
	}
}

func TestClassifyRegisteredWinsOverBuiltin(t *testing.T) {
	defer resetRegistry()

	err := NotFoundf("x")
	Register(err, KindConflict)

	if got := Classify(err); got != KindConflict {
		t.Fatalf("Classify(registered) = %d, want %d", got, KindConflict)
	}
}

func TestIgnoreDefaultBoth(t *testing.T) {
	defer resetRegistry()

	e := New("ignored")
	Ignore(e)

	opt, ok := IgnoreOf(e)
	if !ok {
		t.Fatalf("IgnoreOf(e) ok = false, want true")
	}
	if opt&IgnoreAsSuccess == 0 || opt&IgnoreSilenceLog == 0 {
		t.Fatalf("IgnoreOf(e) = %d, want both bits set", opt)
	}
}

func TestIgnoreSilenceOnly(t *testing.T) {
	defer resetRegistry()

	e := New("ignored")
	Ignore(e, IgnoreSilenceLog)

	opt, ok := IgnoreOf(e)
	if !ok {
		t.Fatalf("IgnoreOf(e) ok = false, want true")
	}
	if opt&IgnoreSilenceLog == 0 {
		t.Fatalf("IgnoreOf(e) missing IgnoreSilenceLog, got %d", opt)
	}
	if opt&IgnoreAsSuccess != 0 {
		t.Fatalf("IgnoreOf(e) unexpectedly has IgnoreAsSuccess, got %d", opt)
	}
}

func TestIgnoreOfUnregistered(t *testing.T) {
	defer resetRegistry()

	if _, ok := IgnoreOf(New("x")); ok {
		t.Fatalf("IgnoreOf(unregistered) ok = true, want false")
	}
}
