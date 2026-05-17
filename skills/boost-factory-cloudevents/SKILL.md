---
name: boost-factory-cloudevents
description: "Use when constructing a CloudEvents HTTP receiver or sender via github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2. Covers NewHTTP for handler/server-style HTTP CloudEvents reception. Triggers on imports under factory/contrib/cloudevents/sdk-go/, on questions about CloudEvents HTTP transport or NewHTTP construction. For Pub/Sub / NATS / Kafka CloudEvents flow, see boost-bootstrap-adapter-{pubsub,nats,kafka}."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-bootstrap-function` (handler typing). For pub/sub-based CloudEvents flows → `boost-bootstrap-adapter-pubsub` / `-nats` / `-kafka`.

```go
import cefact "github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2"

c := cefact.NewHTTP(ctx, handler)
// c is a cloudevents.Client wired with HTTP transport
```

Use this when your function is invoked via HTTP CloudEvents (Knative push, Cloud Run Eventarc HTTP trigger, GitHub webhook bridge) instead of broker-pull. The handler signature follows `boost-bootstrap-function`'s `Handler[T]` rule (input value, output pointer).

## Red flags

| Red flag | Fix |
|---|---|
| `cloudevents.NewClientHTTP(...)` from upstream SDK directly | `cefact.NewHTTP(ctx, handler)` |
| Bypassing the function middleware chain when running over HTTP | Wire `function.Wrapper[*cloudevents.Event]` over the handler before passing to `NewHTTP` |
