---
name: boost-bootstrap-adapter-nats
description: "Use when writing or reviewing a Go event-driven service that subscribes to NATS via github.com/xgodev/boost/bootstrap/function/adapter/contrib/nats-io/nats.go/v1. Covers subscriber wiring, queue groups, and the same ctx-loss issue documented for the Pub/Sub adapter (helper.go and subscriber.go hard-code context.Background, requiring the workaround pattern for graceful shutdown). Triggers on imports under bootstrap/function/adapter/contrib/nats-io/, on questions about NATS subscribers in a boost function, or on signal handling for NATS workers."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-bootstrap-function` — handler typing rule.
- `boost-bootstrap-middleware` — recovery/logger/publisher chain.
- `boost-extra-middleware` — `NewAnyErrorWrapper` for the workaround.
- `boost-bootstrap-adapter-pubsub` — same shape, full workaround pattern documented there.

## Canonical (prototype / dev)

```go
import (
    anats "github.com/xgodev/boost/bootstrap/function/adapter/contrib/nats-io/nats.go/v1"
    "github.com/xgodev/boost/bootstrap/function"
)

fn, _ := function.New[*cloudevents.Event](rec, lmi, pmi)
fn.Run(ctx, handle, anats.New[*cloudevents.Event](conn))
```

Subscriptions are configured via `boost.bootstrap.function.adapter.nats.subjects` (comma-separated) and `boost.bootstrap.function.adapter.nats.queueGroup`. Override at deploy via `BOOST_BOOTSTRAP_FUNCTION_ADAPTER_NATS_*`.

## Production caveat — same ctx-loss as Pub/Sub

`bootstrap/function/adapter/contrib/nats-io/nats.go/v1/helper.go:44/46` and `subscriber.go:62` hard-code `context.Background()`. SIGTERM does not gracefully drain.

Apply the **same workaround pattern** documented in `boost-bootstrap-adapter-pubsub`: bypass `fn.Run`, build the chain via `extra/middleware.NewAnyErrorWrapper`, drive `anats.NewSubscriber` with a signal-aware ctx, and add the `// TODO(boost-upstream):` annotation naming the offending file.

## Queue groups

NATS queue groups distribute messages across N subscribers (load balancing). Configure via `boost.bootstrap.function.adapter.nats.queueGroup`. Without a queue group, every subscriber gets every message (broadcast).

## Red flags

| Red flag | Fix |
|---|---|
| `nats.Conn.Subscribe(...)` directly from the upstream SDK | Use `anats.NewSubscriber(...).Subscribe(ctx)` or `function.New + fn.Run` |
| Bypass of `fn.Run` without `// TODO(boost-upstream):` naming `helper.go:44`/`subscriber.go:62` | Add the comment, OR accept ungraceful shutdown |
| Multiple NATS connections per process | Construct one `*nats.Conn` at startup, share |
| Config tunables read via `os.Getenv` | Use `BOOST_BOOTSTRAP_FUNCTION_ADAPTER_NATS_*` overrides |
