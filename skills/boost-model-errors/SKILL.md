---
name: boost-model-errors
description: "Use when creating, wrapping, or matching errors in a Go service that imports github.com/xgodev/boost/model/errors. Covers the typed error catalog (BadRequest, NotFound, Conflict, Forbidden, Internal, NotValid, etc.), how Echo's error_handler plugin and the function publisher deadletter middleware match on these types, why fmt.Errorf(%w) defeats both, and how to register custom (non-boost) errors so they map to HTTP/gRPC codes or get ignored. Triggers on imports of github.com/xgodev/boost/model/errors, on questions about error wrapping in boost, on echo.NewHTTPError uses in a boost handler, on Wrap / NotValidf / Internalf / NewBadRequest naming, or on errors.Register / RegisterMatch / Classify / Ignore / Kind, mapping a custom error to a status code, or ignoring an error."
license: MIT
metadata:
  author: jpfaria
  version: "0.3.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

## Typed errors

```go
import bootsterrors "github.com/xgodev/boost/model/errors"

return bootsterrors.NewBadRequest(err, "invalid payload")        // wraps upstream err
return bootsterrors.BadRequestf("field %q is required", "id")    // standalone
return bootsterrors.NewNotFound(err, "order not found")
return bootsterrors.NewConflict(err, "duplicate order id")
return bootsterrors.NewForbidden(err, "missing permission")
return bootsterrors.NewInternal(err, "downstream call failed")
return bootsterrors.NotValidf("invalid event data")              // for function deadletter
return bootsterrors.Wrap(err, bootsterrors.NotValidf("..."))     // wrap + classify
```

## Type-name matching is what makes them useful

Two boost subsystems pattern-match on the unwrapped error type name:

| Matcher | Where | Routes by |
|---|---|---|
| Echo `error_handler` plugin | HTTP responses (see `boost-factory-echo`) | `*errors.NotFound` → 404, `*errors.BadRequest` → 400, `*errors.Conflict` → 409, `*errors.Forbidden` → 403, `*errors.Internal` → 500, `*errors.NotValid` → 422 |
| Function `publisher` middleware (deadletter mode) | Event handlers (see `boost-bootstrap-middleware`) | `NotValid` → `notvalid` deadletter topic; `Internal` → retry / alerting; etc. |

`fmt.Errorf("%w", err)` defeats both: the matchers walk boost's `Cause()` (the `causer` interface), and a stdlib `*fmt.wrapError` is **not** a `causer`, so the boost type underneath stays invisible. Never wrap a boost error with `fmt.Errorf` — see **Wrapping & propagation** below.

The HTTP (Echo + function/CloudEvents) and gRPC error handlers resolve the status
via `bootsterrors.Classify(err) Kind` — registered custom errors first, then the
built-in `Is*` catalog. The `Kind → HTTP status` table lives in
`model/restresponse` (`HTTPStatusFor`) and is shared by Echo and the function
adapter; the `Kind → gRPC code` table lives in the gRPC `server` package.

## Registering custom errors (map your own error → code, or ignore)

For an application error that is **not** a boost type, register it once at boot
(before serving) against a semantic `Kind`. It then resolves to the right HTTP
status and gRPC code across all three transports — no per-transport wiring.

```go
import bootsterrors "github.com/xgodev/boost/model/errors"

// "XptoError behaves like NotFound" → HTTP 404 / gRPC NotFound, everywhere.
bootsterrors.Register(ErrXpto, bootsterrors.KindNotFound)         // matches via errors.Is

bootsterrors.RegisterMatch(func(err error) bool {                // matches by type
    _, ok := bootsterrors.Cause(err).(*MyTypedErr)               // Cause, like Is* does
    return ok
}, bootsterrors.KindConflict)

// Ignore: no opts == treat as success (HTTP 200 / gRPC OK) AND silence the log.
bootsterrors.Ignore(ErrNoise)
bootsterrors.Ignore(ErrAudit, bootsterrors.IgnoreSilenceLog)     // status normal, just no log
bootsterrors.Ignore(ErrExpected, bootsterrors.IgnoreAsSuccess)   // 200/OK, still logged
bootsterrors.Ignore(ErrBoth, bootsterrors.IgnoreAsSuccess, bootsterrors.IgnoreSilenceLog) // explicit both
```

Opts are OR-combined; passing none is equivalent to passing both. `Ignore` matches
via `errors.Is`; `IgnoreMatch(func(error) bool, ...IgnoreOption)` matches by predicate.

`Kind` values mirror the catalog: `KindNotFound`, `KindBadRequest`, `KindNotValid`,
`KindConflict`, `KindAlreadyExists`, `KindForbidden`, `KindUnauthorized`,
`KindServiceUnavailable`, `KindNotImplemented`, `KindNotProvisioned`,
`KindNotSupported`, `KindNotAssigned`, `KindMethodNotAllowed`,
`KindTooManyRequests`, `KindTimeout`, `KindInternal` (default).

Precedence in `Classify`: registered matchers (registration order, first wins) →
built-in `Is*` → `KindInternal`. The `match` predicate runs while the registry
lock is held — it must not call back into the registry (`Register`/`Classify`/
`Ignore`/…) or it self-deadlocks.

## Wrapping & propagation — never `fmt.Errorf`

boost **classification** (`Classify` / `Is*`) walks `Cause()` — single-level — **not** stdlib `Unwrap()`. The consequence that still bites:

- `fmt.Errorf("ctx: %w", boostErr)` makes `IsServiceUnavailable(...)` / `Classify(...)` return the WRONG kind — `Cause()` is single-level and a `*fmt.wrapError` isn't a `causer`, so the boost type underneath is invisible → `KindInternal` / 500 leaks out the edge.

Note: boost typed errors **do** implement stdlib `Unwrap()`, so `errors.Is` / `errors.As` traverse them — `errors.Is(boostErr, context.Canceled)` works, and propagating or `Annotate`-ing a boost error keeps stdlib sentinel matching. Only boost's own `Cause()`-based **classification** is single-level, which is why `fmt.Errorf("%w")` still defeats the edge's kind mapping even though `errors.Is` would see through it.

So where `err` is already a boost error (app / use-case / propagating layers), **don't wrap — propagate**:

```go
if err != nil {
    return Result{}, err   // already boost-typed: keeps its kind AND stdlib errors.Is(_, sentinel)
}
```

Add classification only at the boundary where the cause is **non-boost** — with a boost constructor, never `fmt.Errorf`:

| You have | Use |
|---|---|
| an already-boost `err` to bubble up | `return X, err` (propagate) |
| a real non-boost cause (json/driver/transport `err`) | `errors.NewServiceUnavailable(err, msg)` / `NewInternal(err, msg)` |
| a synthetic cause (HTTP status, "not found in cart") | `errors.ServiceUnavailablef("…%d…", code)` (no separate cause) |
| boot/config parse failure | `errors.NewInternal(err, "config: …")` |

`fmt.Sprintf(...)` for building a message string stays fine — only `fmt.Errorf` is banned.

## GraphQL transport — gqlgen `ErrorPresenter`

The Echo / gRPC / function edges map boost errors to a code automatically. GraphQL (gqlgen) returns **HTTP 200** with the error inside the `errors[]` array, so the Echo `error_handler` never sees a resolver error. Wire the equivalent at the gqlgen layer — the 4th transport over the same catalog:

```go
srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
    gqlErr := graphql.DefaultErrorPresenter(ctx, e)
    if gqlErr.Extensions == nil {
        gqlErr.Extensions = map[string]any{}
    }
    // Classify()/Is* use single-level Cause() and don't cross fmt.Errorf("%w");
    // walk the stdlib chain (boost errors implement Unwrap) and match each node.
    for err := e; err != nil; err = errors.Unwrap(err) {
        switch {
        case bootsterrors.IsNotFound(err):
            gqlErr.Extensions["code"] = "NOT_FOUND"
        case bootsterrors.IsServiceUnavailable(err):
            gqlErr.Extensions["code"] = "UPSTREAM_UNAVAILABLE"
        case bootsterrors.IsForbidden(err), bootsterrors.IsConflict(err):
            gqlErr.Extensions["code"] = "FORBIDDEN"
        case bootsterrors.IsBadRequest(err), bootsterrors.IsNotValid(err):
            gqlErr.Extensions["code"] = "BAD_USER_INPUT"
        default:
            continue
        }
        return gqlErr
    }
    gqlErr.Extensions["code"] = "INTERNAL_SERVER_ERROR"
    gqlErr.Message = "internal server error" // never leak internals
    return gqlErr
})
```

Same boost catalog as the other transports; only the public code strings differ. Resolvers/use cases keep returning boost typed errors — the presenter is the only GraphQL-specific piece.

## Red flags

| Red flag | Fix |
|---|---|
| `fmt.Errorf("%w", err)` for an error that flows through Echo or function middleware | `bootsterrors.Wrap(err, bootsterrors.<Type>(...))` |
| `fmt.Errorf("ctx: %w", boostErr)` to add a breadcrumb in app/use-case code | Propagate: `return X, err`. The boost kind + stdlib `errors.Is` survive; a breadcrumb isn't worth losing classification. |
| `fmt.Errorf("%w", boostErr)` then expecting `Classify` / `Is*` to still classify it | They won't — classification uses single-level `Cause()`, not stdlib unwrap, so the boost kind underneath is invisible. Propagate the boost error, or re-classify with a boost constructor. (stdlib `errors.Is` *would* see through, but the edge's kind→code mapping won't.) |
| `echo.NewHTTPError(404, "...")` in a handler | `bootsterrors.NewNotFound(err, "...")` |
| Returning a raw upstream error to a handler caller | Wrap with the right `bootsterrors.New<Type>` so the matcher can route it |
| Inventing a custom error struct for things `model/errors` already covers | Use the existing type — extending the catalog needs an upstream PR, not a local workaround |
| A genuinely app-specific error returning 500 because the handler doesn't know it | `bootsterrors.Register(err, Kind…)` / `RegisterMatch(...)` at boot so it maps to the right code in HTTP and gRPC |
| Editing the Echo/gRPC switch by hand to add a case for your error | Register it instead — the switch is now a shared `Classify` + `Kind` table |
