---
name: boost-extra-middleware
description: "Use when composing a generic middleware chain outside the function-bootstrap presets — typically in the documented production workaround for the Pub/Sub adapter ctx-loss. Covers extra/middleware.NewAnyErrorWrapper[T] and how it composes the same recovery/logger/publisher middlewares that fn.Run wires. Triggers on imports of github.com/xgodev/boost/extra/middleware, on questions about NewAnyErrorWrapper, AnyErrorMiddleware composition, or building a middleware chain by hand."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-bootstrap-middleware` — the middlewares this skill composes.
- `boost-bootstrap-adapter-pubsub` — the canonical reason this skill is reached for (production workaround).

## What it is

`extra/middleware.NewAnyErrorWrapper[T](ctx, name, ...mws)` is the generic chain composer. It returns a wrapper that you can apply to any handler typed compatibly. `fn.Run` uses it internally; you reach for it directly only when you can't use `fn.Run` (e.g., the documented Pub/Sub adapter workaround).

## Usage

```go
import (
    "github.com/xgodev/boost/extra/middleware"
    "github.com/xgodev/boost/bootstrap/function"
)

wrp := middleware.NewAnyErrorWrapper[*cloudevents.Event](
    ctx, "bootstrap", rec, lmi, pmi,
)

wrappedHandler := function.Wrapper[*cloudevents.Event](wrp, handle)

// then drive wrappedHandler manually (e.g., apubsub.NewSubscriber)
```

The middleware order arguments here mirror the order semantics from `boost-bootstrap-middleware`: outermost-to-innermost is `rec, lmi, pmi`.

## When to use

- The Pub/Sub / NATS / Kafka adapter workaround for ctx-loss (see `boost-bootstrap-adapter-pubsub`).
- Custom drivers / adapters where you need the same chain semantics outside `fn.Run`.

**Don't use** when the canonical `function.New + fn.Run` works — `fn.Run` already calls this internally with the right arguments.

## Red flags

| Red flag | Fix |
|---|---|
| Used in place of `fn.Run` without a `// TODO(boost-upstream):` comment explaining why | Either add the comment + tracking issue, or use `fn.Run` |
| Different `T` parameter than the rest of the chain | All on `*cloudevents.Event` |
| Custom middleware order that breaks the recovery-outermost rule | Mirror `boost-bootstrap-middleware`'s canonical order |
