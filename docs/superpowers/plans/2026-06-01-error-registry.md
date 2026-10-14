# Error Registry Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Permitir que a aplicação registre erros customizados, associando-os a um semântico boost (`Kind`) que resolve código HTTP e gRPC nos 3 transportes, ou marcando-os para ignorar.

**Architecture:** Um registry global em `model/errors` mapeia matchers (sentinela via `errors.Is`, ou predicado via `errors.As`) para um `Kind` semântico — sem importar `net/http` nem `grpc/codes`. `Classify(err) Kind` resolve com precedência: registrados → `Is*` embutido → `KindInternal`. As tabelas `Kind → código` ficam na borda de cada transporte (HTTP compartilhada em `restresponse`, gRPC no pacote server). Ignore é configurável por erro (`AsSuccess` e/ou `SilenceLog`).

**Tech Stack:** Go 1.26, testify/suite, labstack/echo v4, google.golang.org/grpc, cloudevents/sdk-go v2.

---

## File Structure

Novos:
- `model/errors/kind.go` — enum `Kind` + `builtinKind(err) Kind` (mapeia `Is*` → `Kind`).
- `model/errors/registry.go` — registry global, `Register`/`RegisterMatch`/`Classify`, `Ignore`/`IgnoreMatch`/`IgnoreOf`, `IgnoreOption`.
- `model/errors/kind_test.go`, `model/errors/registry_test.go`.

Alterados:
- `model/restresponse/status_code.go` (novo arquivo no pacote `response`) — `HTTPStatusFor(errors.Kind) int`.
- `factory/contrib/labstack/echo/v4/error_handler.go` — `ErrorStatusCode` via `Classify`; honrar `IgnoreAsSuccess`.
- `factory/contrib/google.golang.org/grpc/v1/server/error.go` — `Error` via `Classify`; honrar `IgnoreAsSuccess`.
- `bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http/run.go` — `ErrorStatusCode` via `Classify`; honrar `IgnoreAsSuccess`.
- `factory/contrib/labstack/echo/v4/plugins/local/wrapper/log/log.go` — `IgnoreSilenceLog` baixa nível.
- `factory/contrib/google.golang.org/grpc/v1/server/plugins/local/wrapper/log/log.go` — `IgnoreSilenceLog` omite campo `error` e baixa nível.
- Skills + READMEs + `.claude-plugin/plugin.json` (bump).

**Comportamento preservado:** `validator.ValidationErrors` continua tratado em cada transporte (HTTP 422 / gRPC InvalidArgument), fora do `Classify` (model/errors não importa validator). Todos os `Is*` ficam inalterados.

**Gap fechado de brinde:** `KindTimeout` (HTTP 408 / gRPC DeadlineExceeded) e `KindTooManyRequests` (HTTP 429 / gRPC ResourceExhausted) — antes caíam em 500/Internal.

---

## Task 1: Enum `Kind` + classificação embutida

**Files:**
- Create: `model/errors/kind.go`
- Test: `model/errors/kind_test.go`

- [ ] **Step 1: Write the failing test**

```go
package errors

import "testing"

func TestBuiltinKind(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want Kind
	}{
		{"notfound", NotFoundf("x"), KindNotFound},
		{"badrequest", BadRequestf("x"), KindBadRequest},
		{"notvalid", NotValidf("x"), KindNotValid},
		{"conflict", Conflictf("x"), KindConflict},
		{"alreadyexists", AlreadyExistsf("x"), KindAlreadyExists},
		{"forbidden", Forbiddenf("x"), KindForbidden},
		{"unauthorized", Unauthorizedf("x"), KindUnauthorized},
		{"serviceunavailable", ServiceUnavailablef("x"), KindServiceUnavailable},
		{"notimplemented", NotImplementedf("x"), KindNotImplemented},
		{"notprovisioned", NotProvisionedf("x"), KindNotProvisioned},
		{"notsupported", NotSupportedf("x"), KindNotSupported},
		{"notassigned", NotAssignedf("x"), KindNotAssigned},
		{"methodnotallowed", MethodNotAllowedf("x"), KindMethodNotAllowed},
		{"toomanyrequests", TooManyRequestsf("x"), KindTooManyRequests},
		{"timeout", Timeoutf("x"), KindTimeout},
		{"internal", Internalf("x"), KindInternal},
		{"plain", New("x"), KindInternal},
	}
	for _, c := range cases {
		if got := builtinKind(c.err); got != c.want {
			t.Errorf("%s: builtinKind = %v, want %v", c.name, got, c.want)
		}
	}
}
```

> Antes de escrever o teste, confirme os nomes exatos dos construtores `*f`
> em `model/errors/*.go` (ex.: `grep -rn "^func .*f(format" model/errors`).
> Ajuste o teste se algum construtor tiver assinatura diferente.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./model/errors/ -run TestBuiltinKind -v`
Expected: FAIL — `undefined: Kind` / `undefined: builtinKind`.

- [ ] **Step 3: Write minimal implementation**

```go
package errors

// Kind is a transport-agnostic semantic classification of an error.
// Each Kind mirrors a category in the boost error catalog and is resolved
// to a concrete HTTP status / gRPC code at the transport boundary.
type Kind int

const (
	// KindInternal is the default/fallback classification.
	KindInternal Kind = iota
	KindNotFound
	KindBadRequest
	KindNotValid
	KindConflict
	KindAlreadyExists
	KindForbidden
	KindUnauthorized
	KindServiceUnavailable
	KindNotImplemented
	KindNotProvisioned
	KindNotSupported
	KindNotAssigned
	KindMethodNotAllowed
	KindTooManyRequests
	KindTimeout
)

// builtinKind classifies err using the built-in Is* catalog checks.
// Returns KindInternal when no boost type matches.
func builtinKind(err error) Kind {
	switch {
	case IsNotFound(err):
		return KindNotFound
	case IsMethodNotAllowed(err):
		return KindMethodNotAllowed
	case IsNotValid(err):
		return KindNotValid
	case IsBadRequest(err):
		return KindBadRequest
	case IsServiceUnavailable(err):
		return KindServiceUnavailable
	case IsConflict(err):
		return KindConflict
	case IsAlreadyExists(err):
		return KindAlreadyExists
	case IsNotImplemented(err):
		return KindNotImplemented
	case IsNotProvisioned(err):
		return KindNotProvisioned
	case IsUnauthorized(err):
		return KindUnauthorized
	case IsForbidden(err):
		return KindForbidden
	case IsNotSupported(err):
		return KindNotSupported
	case IsNotAssigned(err):
		return KindNotAssigned
	case IsTooManyRequests(err):
		return KindTooManyRequests
	case IsTimeout(err):
		return KindTimeout
	default:
		return KindInternal
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./model/errors/ -run TestBuiltinKind -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add model/errors/kind.go model/errors/kind_test.go
git commit -m "feat(errors): add Kind enum and builtin classification"
```

---

## Task 2: Registry de classificação (`Register`/`RegisterMatch`/`Classify`)

**Files:**
- Create: `model/errors/registry.go`
- Test: `model/errors/registry_test.go`

- [ ] **Step 1: Write the failing test**

```go
package errors

import (
	stderrors "errors"
	"testing"
)

type xptoError struct{ msg string }

func (e *xptoError) Error() string { return e.msg }

func TestClassifyRegisteredSentinel(t *testing.T) {
	defer resetRegistry()
	sentinel := New("sentinel boom")
	Register(sentinel, KindNotFound)

	if got := Classify(sentinel); got != KindNotFound {
		t.Fatalf("Classify = %v, want KindNotFound", got)
	}
}

func TestClassifyRegisteredMatch(t *testing.T) {
	defer resetRegistry()
	RegisterMatch(func(err error) bool {
		var x *xptoError
		return stderrors.As(err, &x)
	}, KindNotFound)

	if got := Classify(&xptoError{"boom"}); got != KindNotFound {
		t.Fatalf("Classify = %v, want KindNotFound", got)
	}
}

func TestClassifyFallsBackToBuiltin(t *testing.T) {
	defer resetRegistry()
	if got := Classify(NotFoundf("x")); got != KindNotFound {
		t.Fatalf("Classify = %v, want KindNotFound", got)
	}
}

func TestClassifyDefaultInternal(t *testing.T) {
	defer resetRegistry()
	if got := Classify(New("anything")); got != KindInternal {
		t.Fatalf("Classify = %v, want KindInternal", got)
	}
}

func TestClassifyRegisteredWinsOverBuiltin(t *testing.T) {
	defer resetRegistry()
	// A boost NotFound error, re-registered as Conflict, must classify as Conflict.
	e := NotFoundf("x")
	Register(e, KindConflict)
	if got := Classify(e); got != KindConflict {
		t.Fatalf("Classify = %v, want KindConflict", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./model/errors/ -run TestClassify -v`
Expected: FAIL — `undefined: Register` / `resetRegistry`.

- [ ] **Step 3: Write minimal implementation**

```go
package errors

import "sync"

type matcher struct {
	match func(error) bool
	kind  Kind
}

var (
	registryMu sync.RWMutex
	matchers   []matcher
)

// Register associates target with kind. An error err is matched when
// errors.Is(err, target) is true. Intended to be called during boot,
// before serving requests.
func Register(target error, kind Kind) {
	RegisterMatch(func(err error) bool { return Is(err, target) }, kind)
}

// RegisterMatch associates a custom predicate with kind. Use this for
// type-based matching (errors.As) or any arbitrary condition.
func RegisterMatch(match func(error) bool, kind Kind) {
	if match == nil {
		return
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	matchers = append(matchers, matcher{match: match, kind: kind})
}

// Classify resolves err to a Kind. Precedence: registered matchers (in
// registration order, first match wins) -> built-in Is* catalog ->
// KindInternal.
func Classify(err error) Kind {
	if err == nil {
		return KindInternal
	}
	registryMu.RLock()
	for _, m := range matchers {
		if m.match(err) {
			registryMu.RUnlock()
			return m.kind
		}
	}
	registryMu.RUnlock()
	return builtinKind(err)
}

// resetRegistry clears all registrations. Test-only helper.
func resetRegistry() {
	registryMu.Lock()
	defer registryMu.Unlock()
	matchers = nil
	ignores = nil
}
```

> `resetRegistry` referencia `ignores`, definido na Task 3. Se rodar a Task 2
> isolada antes da 3, troque temporariamente o corpo para `matchers = nil` e
> reponha `ignores = nil` na Task 3.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./model/errors/ -run TestClassify -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add model/errors/registry.go model/errors/registry_test.go
git commit -m "feat(errors): add custom error classification registry"
```

---

## Task 3: Ignore configurável (`Ignore`/`IgnoreMatch`/`IgnoreOf`)

**Files:**
- Modify: `model/errors/registry.go`
- Test: `model/errors/registry_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestIgnoreDefaultBoth(t *testing.T) {
	defer resetRegistry()
	e := New("noise")
	Ignore(e) // no opts == AsSuccess + SilenceLog

	pol, ok := IgnoreOf(e)
	if !ok {
		t.Fatal("IgnoreOf: ok = false, want true")
	}
	if pol&IgnoreAsSuccess == 0 || pol&IgnoreSilenceLog == 0 {
		t.Fatalf("IgnoreOf = %b, want both bits set", pol)
	}
}

func TestIgnoreSilenceOnly(t *testing.T) {
	defer resetRegistry()
	e := New("quiet")
	Ignore(e, IgnoreSilenceLog)

	pol, ok := IgnoreOf(e)
	if !ok || pol&IgnoreSilenceLog == 0 || pol&IgnoreAsSuccess != 0 {
		t.Fatalf("IgnoreOf = %b ok=%v, want SilenceLog only", pol, ok)
	}
}

func TestIgnoreOfUnregistered(t *testing.T) {
	defer resetRegistry()
	if _, ok := IgnoreOf(New("x")); ok {
		t.Fatal("IgnoreOf: ok = true, want false")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./model/errors/ -run TestIgnore -v`
Expected: FAIL — `undefined: Ignore` / `IgnoreOf` / `IgnoreAsSuccess`.

- [ ] **Step 3: Write minimal implementation**

Append to `model/errors/registry.go`:

```go
// IgnoreOption configures how an ignored error is handled. Combine bits.
type IgnoreOption int

const (
	// IgnoreAsSuccess makes the transport respond with success
	// (HTTP 200/204, gRPC codes.OK) instead of an error.
	IgnoreAsSuccess IgnoreOption = 1 << iota
	// IgnoreSilenceLog suppresses (or lowers) logging of the error.
	IgnoreSilenceLog
)

type ignoreRule struct {
	match  func(error) bool
	policy IgnoreOption
}

var ignores []ignoreRule

// Ignore marks errors matching target (via errors.Is) to be ignored.
// With no options, the error is fully ignored (AsSuccess + SilenceLog).
func Ignore(target error, opts ...IgnoreOption) {
	IgnoreMatch(func(err error) bool { return Is(err, target) }, opts...)
}

// IgnoreMatch marks errors matching the predicate to be ignored.
// With no options, the error is fully ignored (AsSuccess + SilenceLog).
func IgnoreMatch(match func(error) bool, opts ...IgnoreOption) {
	if match == nil {
		return
	}
	var policy IgnoreOption
	if len(opts) == 0 {
		policy = IgnoreAsSuccess | IgnoreSilenceLog
	} else {
		for _, o := range opts {
			policy |= o
		}
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	ignores = append(ignores, ignoreRule{match: match, policy: policy})
}

// IgnoreOf returns the ignore policy for err, if any was registered.
func IgnoreOf(err error) (IgnoreOption, bool) {
	if err == nil {
		return 0, false
	}
	registryMu.RLock()
	defer registryMu.RUnlock()
	for _, r := range ignores {
		if r.match(err) {
			return r.policy, true
		}
	}
	return 0, false
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./model/errors/ -run TestIgnore -v`
Expected: PASS.

- [ ] **Step 5: Run the whole package + commit**

Run: `go test ./model/errors/...`
Expected: PASS (todos os testes, inclusive os `Is*` existentes).

```bash
git add model/errors/registry.go model/errors/registry_test.go
git commit -m "feat(errors): add configurable ignore policy registry"
```

---

## Task 4: Tabela HTTP compartilhada (`restresponse.HTTPStatusFor`)

**Files:**
- Create: `model/restresponse/status_code.go`
- Test: `model/restresponse/status_code_test.go`

> Verifique antes que `model/restresponse` ainda não importe `model/errors`
> e que `model/errors` não importe `restresponse` (sem ciclo). `model/errors`
> só importa stdlib, então é seguro.

- [ ] **Step 1: Write the failing test**

```go
package response

import (
	"net/http"
	"testing"

	"github.com/xgodev/boost/model/errors"
)

func TestHTTPStatusFor(t *testing.T) {
	cases := map[errors.Kind]int{
		errors.KindNotFound:           http.StatusNotFound,
		errors.KindMethodNotAllowed:   http.StatusMethodNotAllowed,
		errors.KindNotValid:           http.StatusBadRequest,
		errors.KindBadRequest:         http.StatusBadRequest,
		errors.KindServiceUnavailable: http.StatusServiceUnavailable,
		errors.KindConflict:           http.StatusConflict,
		errors.KindAlreadyExists:      http.StatusConflict,
		errors.KindNotImplemented:     http.StatusNotImplemented,
		errors.KindNotProvisioned:     http.StatusNotImplemented,
		errors.KindUnauthorized:       http.StatusUnauthorized,
		errors.KindForbidden:          http.StatusForbidden,
		errors.KindNotSupported:       http.StatusUnprocessableEntity,
		errors.KindNotAssigned:        http.StatusUnprocessableEntity,
		errors.KindTooManyRequests:    http.StatusTooManyRequests,
		errors.KindTimeout:            http.StatusRequestTimeout,
		errors.KindInternal:           http.StatusInternalServerError,
	}
	for kind, want := range cases {
		if got := HTTPStatusFor(kind); got != want {
			t.Errorf("HTTPStatusFor(%v) = %d, want %d", kind, got, want)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./model/restresponse/ -run TestHTTPStatusFor -v`
Expected: FAIL — `undefined: HTTPStatusFor`.

- [ ] **Step 3: Write minimal implementation**

```go
package response

import (
	"net/http"

	"github.com/xgodev/boost/model/errors"
)

// HTTPStatusFor maps a boost error Kind to its HTTP status code.
func HTTPStatusFor(kind errors.Kind) int {
	switch kind {
	case errors.KindNotFound:
		return http.StatusNotFound
	case errors.KindMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case errors.KindNotValid, errors.KindBadRequest:
		return http.StatusBadRequest
	case errors.KindServiceUnavailable:
		return http.StatusServiceUnavailable
	case errors.KindConflict, errors.KindAlreadyExists:
		return http.StatusConflict
	case errors.KindNotImplemented, errors.KindNotProvisioned:
		return http.StatusNotImplemented
	case errors.KindUnauthorized:
		return http.StatusUnauthorized
	case errors.KindForbidden:
		return http.StatusForbidden
	case errors.KindNotSupported, errors.KindNotAssigned:
		return http.StatusUnprocessableEntity
	case errors.KindTooManyRequests:
		return http.StatusTooManyRequests
	case errors.KindTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./model/restresponse/ -run TestHTTPStatusFor -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add model/restresponse/status_code.go model/restresponse/status_code_test.go
git commit -m "feat(restresponse): add Kind->HTTP status mapping"
```

---

## Task 5: Refactor do error handler do Echo + IgnoreAsSuccess

**Files:**
- Modify: `factory/contrib/labstack/echo/v4/error_handler.go`
- Test: `factory/contrib/labstack/echo/v4/error_handler_test.go`

- [ ] **Step 1: Write the failing test**

```go
package echo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/model/errors"
)

type customNotFound struct{ msg string }

func (c *customNotFound) Error() string { return c.msg }

func TestErrorStatusCode_RegisteredCustom(t *testing.T) {
	errors.RegisterMatch(func(err error) bool {
		_, ok := err.(*customNotFound)
		return ok
	}, errors.KindNotFound)

	if got := ErrorStatusCode(&customNotFound{"missing"}); got != http.StatusNotFound {
		t.Fatalf("ErrorStatusCode = %d, want 404", got)
	}
}

func TestErrorHandler_IgnoreAsSuccess(t *testing.T) {
	sentinel := errors.New("ignorable")
	errors.Ignore(sentinel, errors.IgnoreAsSuccess)

	rec := httptest.NewRecorder()
	srv := e.New()
	c := srv.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

	errorHandler(sentinel, c, e.MIMEApplicationJSON)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./factory/contrib/labstack/echo/v4/ -run TestErrorStatusCode_RegisteredCustom -v`
Expected: FAIL — comportamento atual retorna 500 para `*customNotFound`.

- [ ] **Step 3: Write minimal implementation**

Substitua o corpo de `errorHandler` e `ErrorStatusCode` em `error_handler.go`
(mantendo imports `net/http`, `strconv`, validator, `response`; adicione o
import do pacote errors caso ainda não exista):

```go
func errorHandler(err error, c e.Context, contentType string) {
	// Ignored-as-success errors short-circuit to a success response.
	if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreAsSuccess != 0 {
		if er := c.NoContent(http.StatusOK); er != nil {
			c.Logger().Error(er)
		}
		return
	}

	var (
		status  int
		message string
	)
	if echoErr, ok := err.(*e.HTTPError); ok {
		status = echoErr.Code
		message = fmt.Sprintf("%v", echoErr.Message)
	} else {
		status = ErrorStatusCode(err)
		message = err.Error()
	}

	var er error
	if c.Request().Method == http.MethodHead {
		er = c.NoContent(status)
	} else {
		switch contentType {
		case e.MIMEApplicationJSON:
			er = c.JSON(status, response.Error{HttpStatusCode: status, ErrorCode: strconv.Itoa(status), Message: message})
		default:
			er = c.String(status, message)
		}
	}
	if er != nil {
		c.Logger().Error(er)
	}
}

// ErrorStatusCode translates err to the respective HTTP status code.
func ErrorStatusCode(err error) int {
	if _, ok := err.(validator.ValidationErrors); ok {
		return http.StatusUnprocessableEntity
	}
	return response.HTTPStatusFor(errors.Classify(err))
}
```

Garanta o import `response "github.com/xgodev/boost/model/restresponse"` (já
existe) e `"github.com/xgodev/boost/model/errors"` (já existe).

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./factory/contrib/labstack/echo/v4/ -run 'TestErrorStatusCode|TestErrorHandler' -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add factory/contrib/labstack/echo/v4/error_handler.go factory/contrib/labstack/echo/v4/error_handler_test.go
git commit -m "refactor(echo): map errors via registry Classify + honor ignore"
```

---

## Task 6: Refactor do error mapping do gRPC + IgnoreAsSuccess

**Files:**
- Modify: `factory/contrib/google.golang.org/grpc/v1/server/error.go`
- Test: `factory/contrib/google.golang.org/grpc/v1/server/error_test.go`

- [ ] **Step 1: Write the failing test**

```go
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

	got := status.Code(Error(&customConflict{"dup"}))
	if got != codes.AlreadyExists {
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./factory/contrib/google.golang.org/grpc/v1/server/ -run TestError_ -v`
Expected: FAIL — custom error mapeia para Internal; ignore não retorna nil.

- [ ] **Step 3: Write minimal implementation**

Substitua `error.go`:

```go
package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/xgodev/boost/model/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error converts a boost/application error to a gRPC status error.
// Errors registered as ignore-as-success return nil.
func Error(err error) error {
	if err == nil {
		return nil
	}
	if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreAsSuccess != 0 {
		return nil
	}
	if _, ok := err.(validator.ValidationErrors); ok {
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}
	return status.Errorf(grpcCodeFor(errors.Classify(err)), "%s", err.Error())
}

// grpcCodeFor maps a boost error Kind to a gRPC code.
func grpcCodeFor(kind errors.Kind) codes.Code {
	switch kind {
	case errors.KindNotFound:
		return codes.NotFound
	case errors.KindNotValid, errors.KindBadRequest:
		return codes.InvalidArgument
	case errors.KindServiceUnavailable:
		return codes.Unavailable
	case errors.KindConflict, errors.KindAlreadyExists:
		return codes.AlreadyExists
	case errors.KindNotImplemented, errors.KindNotProvisioned:
		return codes.Unimplemented
	case errors.KindUnauthorized:
		return codes.Unauthenticated
	case errors.KindForbidden:
		return codes.PermissionDenied
	case errors.KindTooManyRequests:
		return codes.ResourceExhausted
	case errors.KindTimeout:
		return codes.DeadlineExceeded
	default:
		return codes.Internal
	}
}
```

> Nota: `KindMethodNotAllowed`, `KindNotSupported`, `KindNotAssigned` caem em
> `codes.Internal` (preserva o comportamento atual do gRPC, que não tinha case
> para eles).

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./factory/contrib/google.golang.org/grpc/v1/server/ -run TestError_ -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add factory/contrib/google.golang.org/grpc/v1/server/error.go factory/contrib/google.golang.org/grpc/v1/server/error_test.go
git commit -m "refactor(grpc): map errors via registry Classify + honor ignore"
```

---

## Task 7: Refactor do adapter function/CloudEvents + IgnoreAsSuccess

**Files:**
- Modify: `bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http/run.go`
- Test: `bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http/run_test.go`

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http/ -run TestErrorStatusCode -v`
Expected: FAIL — custom error retorna 500.

- [ ] **Step 3: Write minimal implementation**

Substitua a função `ErrorStatusCode` e ajuste `Wrapper` em `run.go`. Adicione
o import `response "github.com/xgodev/boost/model/restresponse"`.

`ErrorStatusCode`:

```go
// ErrorStatusCode translates err to the respective HTTP status code.
func ErrorStatusCode(err error) int {
	if _, ok := err.(validator.ValidationErrors); ok {
		return http.StatusUnprocessableEntity
	}
	return response.HTTPStatusFor(errors.Classify(err))
}
```

No `Wrapper`, trate ignore-as-success antes de montar o erro (bloco `if err != nil`):

```go
		e, err := fn(ctx, event)
		if err != nil {
			if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreAsSuccess != 0 {
				return ce.NewHTTPResult(http.StatusOK, "")
			}
			status := ErrorStatusCode(err)
			return ce.NewHTTPResult(status, err.Error())
		}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http/ -run TestErrorStatusCode -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http/run.go bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http/run_test.go
git commit -m "refactor(function/ce): map errors via registry Classify + honor ignore"
```

---

## Task 8: IgnoreSilenceLog nos middlewares de log (echo + gRPC)

**Files:**
- Modify: `factory/contrib/labstack/echo/v4/plugins/local/wrapper/log/log.go`
- Modify: `factory/contrib/google.golang.org/grpc/v1/server/plugins/local/wrapper/log/log.go`

> Sem teste de unidade dedicado (middlewares dependem de servidor real). A
> verificação é via `go build` + revisão. O efeito: erro com `IgnoreSilenceLog`
> não polui log de erro.

- [ ] **Step 1: Echo — baixar nível quando SilenceLog**

Em `loggerMiddleware`, a `err` da request está em escopo (linha ~119). No
`defer`, antes do `switch level`, force `Debugf` quando o erro for silenciado.
Adicione o import `"github.com/xgodev/boost/model/errors"` e capture `err` numa
variável visível ao defer (renomeie a declaração para antes do defer):

```go
			var err error
			defer func() {
				stop := time.Now()

				reqSize := req.Header.Get(e.HeaderContentLength)
				if reqSize == "" {
					reqSize = "0"
				}

				var method func(format string, args ...interface{})

				if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreSilenceLog != 0 {
					method = logger.Debugf
				} else {
					switch level {
					case "TRACE":
						method = logger.Tracef
					case "INFO":
						method = logger.Infof
					default:
						method = logger.Debugf
					}
				}

				method("%s %s %-7s %s %3d %s %s %13v %s %s",
					c.RealIP(), req.Host, req.Method, req.RequestURI,
					res.Status, reqSize, strconv.FormatInt(res.Size, 10),
					stop.Sub(start).String(), req.Referer(), req.UserAgent(),
				)
			}()

			if err = next(c); err != nil {
				c.Error(err)
			}

			return nil
```

- [ ] **Step 2: gRPC — omitir campo error + baixar nível quando SilenceLog**

Nos dois interceptors (`streamInterceptor` e `unaryInterceptor`), troque o bloco
`if err != nil { logger = logger.WithField("error", err.Error()) }` e a escolha
do método. Adicione o import `"github.com/xgodev/boost/model/errors"`.

`unaryInterceptor` (mesma ideia no `streamInterceptor`):

```go
		silenced := false
		if err != nil {
			if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreSilenceLog != 0 {
				silenced = true
			} else {
				logger = logger.WithField("error", err.Error())
			}
		}

		method := i.m(logger)
		if silenced {
			method = logger.Debugf
		}
		method("unary request received")
		return resp, err
```

(No `streamInterceptor`, a mensagem é `"stream request received"`.)

- [ ] **Step 3: Build**

Run: `go build ./...`
Expected: sem erros.

- [ ] **Step 4: Commit**

```bash
git add factory/contrib/labstack/echo/v4/plugins/local/wrapper/log/log.go factory/contrib/google.golang.org/grpc/v1/server/plugins/local/wrapper/log/log.go
git commit -m "feat(log): honor IgnoreSilenceLog in echo and grpc log middlewares"
```

---

## Task 9: Verificação global (build + vet + test)

**Files:** nenhum (gate).

- [ ] **Step 1: Build, vet, test**

Run:
```bash
gofmt -l model/errors model/restresponse factory bootstrap
go vet ./...
go build ./...
go test ./...
```
Expected: `gofmt -l` vazio; vet sem warning; build ok; testes verdes.

- [ ] **Step 2: Commit (se houve format)**

```bash
git add -A && git commit -m "chore: gofmt" || echo "nada a formatar"
```

---

## Task 10: Docs + skills + bump do plugin (Iron Law #6)

**Files:**
- Modify: `skills/boost-model-errors/SKILL.md` — seção "Registrar erros customizados": `Register`/`RegisterMatch`/`Classify`/`Ignore` + tabela `Kind → HTTP/gRPC`.
- Modify: `skills/boost-factory-echo/SKILL.md` — nota de que o error_handler resolve via registry.
- Modify: `skills/boost-factory-grpc/SKILL.md` — idem para o `Error`.
- Modify: READMEs de `model/errors`, `model/restresponse` (se existirem) — nova API.
- Modify: `.claude-plugin/plugin.json` — bump de versão.

> Antes de editar QUALQUER `SKILL.md`, invoque `superpowers:writing-skills`
> (baseline com subagent primeiro). Antes de criar hook, `hookify:writing-rules`.

- [ ] **Step 1: Atualizar skill boost-model-errors**

Documente o exemplo canônico:

```go
// No boot, antes de servir:
errors.Register(MyXptoError, errors.KindNotFound)        // 404 / NotFound
errors.RegisterMatch(func(e error) bool {                 // por tipo
    var x *MyTypedErr; return errors.As(e, &x)
}, errors.KindConflict)
errors.Ignore(MyNoiseError)                               // 200/OK + sem log
errors.Ignore(MyAuditErr, errors.IgnoreSilenceLog)        // status normal, sem log
```

E a tabela `Kind → HTTP / gRPC` (de `HTTPStatusFor` e `grpcCodeFor`).

- [ ] **Step 2: Atualizar skills echo + grpc + READMEs**

Nota curta: "erros são classificados via `errors.Classify`; registre erros
customizados com `errors.Register`/`RegisterMatch`."

- [ ] **Step 3: Bump do plugin**

Edite `.claude-plugin/plugin.json` subindo a versão (ex.: patch→minor conforme
a convenção do repo).

- [ ] **Step 4: Commit**

```bash
git add skills .claude-plugin/plugin.json model/errors/README.md model/restresponse/README.md 2>/dev/null
git commit -m "docs(errors): document custom error registry in skills and READMEs"
```

---

## Notas de design recapituladas

- `model/errors` permanece agnóstico de transporte: `Kind` é só um inteiro; as
  tabelas `Kind→código` ficam em `restresponse` (HTTP) e no pacote `server` (gRPC).
- Precedência de `Classify`: registrados → `Is*` embutido → `KindInternal`.
- `validator.ValidationErrors` segue tratado em cada transporte (não entra no
  registry), preservando HTTP 422 / gRPC InvalidArgument.
- Ignore é por bits (`IgnoreAsSuccess | IgnoreSilenceLog`); sem opção = ambos.
- `IgnoreSilenceLog` atua nos middlewares de log (echo/gRPC); no adapter CE não
  há log por-evento do erro, então só o `IgnoreAsSuccess` se aplica lá.
- Registro deve ocorrer no boot, antes de servir (registry sob `RWMutex`).
```
