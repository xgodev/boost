---
name: boost-bootstrap-middleware
description: "Use when stacking function middleware in a Go event-driven service that imports github.com/xgodev/boost/bootstrap/function/middleware/{recovery,logger,publisher,ignore_errors}. Covers the canonical recovery -> logger -> publisher chain order, why each layer's position matters (recovery outermost, publisher innermost), and how to wrap errors so the publisher's deadletter mode routes them by type-name. Triggers on imports under bootstrap/function/middleware/, on questions about NewRecovery, NewAnyErrorMiddleware, deadletter routing, or middleware order."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-bootstrap-function` — generic typing rule (`T = *cloudevents.Event`).
- `boost-model-errors` — error types matched by deadletter mode.
- `boost-wrapper-publisher` — provides the publisher consumed by the publisher middleware.

## Canonical chain

```
recovery → logger → publisher
```

```go
import (
    rm "github.com/xgodev/boost/bootstrap/function/middleware/recovery"
    lm "github.com/xgodev/boost/bootstrap/function/middleware/logger"
    pm "github.com/xgodev/boost/bootstrap/function/middleware/publisher"
)

rec := rm.NewRecovery[*cloudevents.Event]()
lmi, _ := lm.NewAnyErrorMiddleware[*cloudevents.Event]()
pmi, _ := pm.NewAnyErrorMiddleware[*cloudevents.Event](pub)

// fn.Run path:
fn, _ := function.New[*cloudevents.Event](rec, lmi, pmi)

// Workaround path (see boost-bootstrap-adapter-pubsub):
wrp := middleware.NewAnyErrorWrapper[*cloudevents.Event](ctx, "bootstrap", rec, lmi, pmi)
```

## Why this order

| Layer | Position | Reason |
|---|---|---|
| `recovery` | Outermost | A panic in the handler must not kill the worker before the logger gets the error |
| `logger` | Middle | Sees both raw handler errors AND post-recovery errors; structured fields propagate |
| `publisher` | Innermost | Fires on the handler's successful result; with deadletter config, also routes errors by type |

Reordering is a footgun: putting `recovery` innermost means a panic short-circuits before the logger; putting `publisher` outermost means it fires before the handler ran.

## Wrapping errors for deadletter routing

The publisher middleware in deadletter mode matches on the unwrapped error type name:

```go
// Routed to a "notvalid" deadletter topic; not retried
return nil, bootsterrors.Wrap(err, bootsterrors.NotValidf("invalid event data"))

// Routed to retry / alerting; transient
return nil, bootsterrors.Wrap(err, bootsterrors.Internalf("downstream call failed"))
```

`fmt.Errorf("%w", err)` defeats this — the matcher cannot recover the type name. Always use `bootsterrors.Wrap` (see `boost-model-errors`).

## `ignore_errors` middleware

Use only when you genuinely want a category of errors to silently ack the message instead of nack/retry. Stack outside the publisher (so the publisher still fires for non-ignored errors), but inside recovery (so panics still recover):

```go
rec := rm.NewRecovery[*cloudevents.Event]()
imi, _ := im.NewAnyErrorMiddleware[*cloudevents.Event]()  // ignore_errors
lmi, _ := lm.NewAnyErrorMiddleware[*cloudevents.Event]()
pmi, _ := pm.NewAnyErrorMiddleware[*cloudevents.Event](pub)

fn, _ := function.New[*cloudevents.Event](rec, imi, lmi, pmi)
```

## Red flags

| Red flag | Fix |
|---|---|
| Chain ordered `publisher → logger → recovery` (or any permutation that violates outermost=recovery) | Reorder to `recovery → logger → publisher` |
| Forgetting `recovery` middleware | Always include it — production functions die otherwise |
| `fmt.Errorf("%w", err)` from a handler returned through this chain | `bootsterrors.Wrap(err, bootsterrors.<Type>(...))` |
| Mixing `T = cloudevents.Event` and `T = *cloudevents.Event` across middlewares | All on `*cloudevents.Event` |
