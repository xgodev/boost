---
name: boost-factory-goka
description: "Use when constructing a Goka emitter (Kafka stream-processing producer) via github.com/xgodev/boost/factory/contrib/lovoo/goka/v1. Covers NewEmitter + variants and the typical use case (publishing to a Goka-managed Kafka stream that downstream Goka processors consume). Triggers on imports under factory/contrib/lovoo/goka/, on questions about Goka emitters or Goka stream processors in a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. For raw Kafka producer/consumer → `boost-factory-kafka`.

```go
import gokafact "github.com/xgodev/boost/factory/contrib/lovoo/goka/v1"

emitter, err := gokafact.NewEmitter(ctx)
if err != nil { log.Fatalf("goka: %v", err) }
defer emitter.Finish()
```

Configure brokers, stream name, codec under `boost.factory.goka.*` (override `BOOST_FACTORY_GOKA_*`).

## Goka vs raw Kafka

Goka adds stream-processing semantics (key/value tables, joins, group processors) on top of Kafka. Use Goka when your service participates in a Goka topology; use raw Kafka (`boost-factory-kafka`) when you just need produce/consume.

The publisher driver `wrapper/publisher/driver/contrib/lovoo/goka/v1` (see `boost-wrapper-publisher`) wraps this emitter for the function publisher middleware path.

## Red flags

| Red flag | Fix |
|---|---|
| `goka.NewEmitter(...)` directly | `gokafact.NewEmitter(ctx)` |
| Brokers via `os.Getenv` | `BOOST_FACTORY_GOKA_*` |
| Forgetting `emitter.Finish()` on shutdown | Add it — flushes pending writes |
