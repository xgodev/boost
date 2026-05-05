---
name: boost-model-errors
description: "Use when creating, wrapping, or matching errors in a Go service that imports github.com/xgodev/boost/model/errors. Covers the typed error catalog (BadRequest, NotFound, Conflict, Forbidden, Internal, NotValid, etc.), how Echo's error_handler plugin and the function publisher deadletter middleware match on these types, and why fmt.Errorf(%w) defeats both. Triggers on imports of github.com/xgodev/boost/model/errors, on questions about error wrapping in boost, on echo.NewHTTPError uses in a boost handler, or on Wrap / NotValidf / Internalf / NewBadRequest naming."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
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

`fmt.Errorf("%w", err)` defeats both because `errors.As(target *NotFound)` can't unwrap an opaque wrapped error of unknown concrete type. Always use `bootsterrors.Wrap`.

## Red flags

| Red flag | Fix |
|---|---|
| `fmt.Errorf("%w", err)` for an error that flows through Echo or function middleware | `bootsterrors.Wrap(err, bootsterrors.<Type>(...))` |
| `echo.NewHTTPError(404, "...")` in a handler | `bootsterrors.NewNotFound(err, "...")` |
| Returning a raw upstream error to a handler caller | Wrap with the right `bootsterrors.New<Type>` so the matcher can route it |
| Inventing a custom error struct for things `model/errors` already covers | Use the existing type — extending the catalog needs an upstream PR, not a local workaround |
