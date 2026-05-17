---
name: boost-factory-gocloud-pubsub
description: "Use when constructing a portable pub/sub topic or subscription via gocloud.dev/pubsub through github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0. Covers NewTopic / NewSubscription + variants and the canonical shape shipped under examples/pubsub/. The URL-driven model lets one binary swap pub/sub provider (gcppubsub://, kafka://, awssns:///, mempubsub://) without code change. Triggers on imports under factory/contrib/gocloud.dev/pubsub/, on questions about gocloud.dev or provider-portable pub/sub in a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`.

## Canonical example (ships with boost)

- `factory/contrib/gocloud.dev/pubsub/v0/examples/pubsub/main.go`

## Construction

```go
import gcppubsub "github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0"

topic, err := gcppubsub.NewTopic(ctx)
if err != nil { log.Fatalf("topic: %v", err) }
defer topic.Shutdown(ctx)

sub, err := gcppubsub.NewSubscription(ctx)
defer sub.Shutdown(ctx)
```

Configure topic/subscription URLs (`gcppubsub://...`, `kafka://...`, `awssns:///...`, `mempubsub://...`) under `boost.factory.gocloud.pubsub.*` (override `BOOST_FACTORY_GOCLOUD_PUBSUB_*`).

## When gocloud.dev vs the native factory?

| Reach for gocloud.dev | Reach for native (`boost-factory-pubsub`, `boost-factory-kafka`, ...) |
|---|---|
| Want to swap providers via URL config without code change | Committed to one provider; want full feature surface |
| Tests use `mempubsub://` for in-process fakes | Production-only path; tests run against real broker |

The native factories expose more provider-specific knobs. gocloud.dev exposes the lowest-common-denominator API.

## Red flags

| Red flag | Fix |
|---|---|
| `pubsub.OpenTopic(ctx, url)` directly | `gcppubsub.NewTopic(ctx)` |
| URL via `os.Getenv` | `BOOST_FACTORY_GOCLOUD_PUBSUB_*` |
| Forgetting `defer topic.Shutdown(ctx)` | Add it |
