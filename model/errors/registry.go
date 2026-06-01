package errors

import "sync"

// matcher associates a predicate with the Kind to return when it matches.
type matcher struct {
	match func(error) bool
	kind  Kind
}

// ignoreRule associates a predicate with the ignore policy to apply.
type ignoreRule struct {
	match  func(error) bool
	policy IgnoreOption
}

// IgnoreOption is a bit flag describing how a matched error should be ignored.
type IgnoreOption int

const (
	// IgnoreAsSuccess treats the matched error as a successful outcome.
	IgnoreAsSuccess IgnoreOption = 1 << iota
	// IgnoreSilenceLog suppresses logging for the matched error.
	IgnoreSilenceLog
)

var (
	mu       sync.RWMutex
	matchers []matcher
	ignores  []ignoreRule
)

// Register registers target so that any error matching it (via Is) is
// classified as kind. For custom predicates see RegisterMatch and its note on
// not calling back into the registry from within a matcher.
func Register(target error, kind Kind) {
	mu.Lock()
	defer mu.Unlock()
	matchers = append(matchers, matcher{
		match: func(err error) bool { return Is(err, target) },
		kind:  kind,
	})
}

// RegisterMatch registers a custom predicate to classify errors as kind. A nil
// match is a no-op. The match function must not call back into this package's
// registry (Register/RegisterMatch/Classify/Ignore/IgnoreMatch/IgnoreOf), as
// matchers run while the registry lock is held.
func RegisterMatch(match func(error) bool, kind Kind) {
	if match == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	matchers = append(matchers, matcher{match: match, kind: kind})
}

// Classify returns the Kind of err. Registered matchers are evaluated in order
// (first match wins); if none match it falls back to builtinKind. A nil err is
// classified as KindInternal.
func Classify(err error) Kind {
	if err == nil {
		return KindInternal
	}

	mu.RLock()
	for _, m := range matchers {
		if m.match(err) {
			kind := m.kind
			mu.RUnlock()
			return kind
		}
	}
	mu.RUnlock()

	return builtinKind(err)
}

// Ignore registers target so that any error matching it (via Is) is ignored
// according to opts. With no opts the full ignore policy
// (IgnoreAsSuccess | IgnoreSilenceLog) is applied. For custom predicates see
// IgnoreMatch and its note on not calling back into the registry from within a
// matcher.
func Ignore(target error, opts ...IgnoreOption) {
	policy := combineOptions(opts)
	mu.Lock()
	defer mu.Unlock()
	ignores = append(ignores, ignoreRule{
		match:  func(err error) bool { return Is(err, target) },
		policy: policy,
	})
}

// IgnoreMatch registers a custom predicate to ignore errors according to opts.
// A nil match is a no-op. With no opts the full ignore policy
// (IgnoreAsSuccess | IgnoreSilenceLog) is applied. The match function must not
// call back into this package's registry
// (Register/RegisterMatch/Classify/Ignore/IgnoreMatch/IgnoreOf), as matchers
// run while the registry lock is held.
func IgnoreMatch(match func(error) bool, opts ...IgnoreOption) {
	if match == nil {
		return
	}
	policy := combineOptions(opts)
	mu.Lock()
	defer mu.Unlock()
	ignores = append(ignores, ignoreRule{match: match, policy: policy})
}

// IgnoreOf returns the ignore policy registered for err and whether one was
// found. Registered rules are evaluated in order (first match wins). A nil err
// returns (0, false).
func IgnoreOf(err error) (IgnoreOption, bool) {
	if err == nil {
		return 0, false
	}
	mu.RLock()
	defer mu.RUnlock()
	for _, r := range ignores {
		if r.match(err) {
			return r.policy, true
		}
	}
	return 0, false
}

// combineOptions ORs the provided opts, defaulting to the full ignore policy
// when none are supplied.
func combineOptions(opts []IgnoreOption) IgnoreOption {
	if len(opts) == 0 {
		return IgnoreAsSuccess | IgnoreSilenceLog
	}
	var policy IgnoreOption
	for _, o := range opts {
		policy |= o
	}
	return policy
}

// resetRegistry clears all registered matchers and ignore rules. Test-only.
func resetRegistry() {
	mu.Lock()
	defer mu.Unlock()
	matchers = nil
	ignores = nil
}
