---
name: boost-factory-buntdb
description: "Use when constructing an embedded buntdb (in-memory key/value store with optional persistence) via github.com/xgodev/boost/factory/contrib/tidwall/buntdb/v1. Covers NewDB + variants and the typical use case (sidecar caches, transient state, dev fixtures). Triggers on imports under factory/contrib/tidwall/buntdb/, on questions about buntdb or embedded KV stores in a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. Layout/trio convention → `boost-maintainer`.

```go
import buntdb "github.com/xgodev/boost/factory/contrib/tidwall/buntdb/v1"

db, err := buntdb.NewDB(ctx)
if err != nil { log.Fatalf("buntdb: %v", err) }
defer db.Close()
```

Configure path, sync mode, auto-shrink under `boost.factory.buntdb.*` (override `BOOST_FACTORY_BUNTDB_*`). Use `:memory:` for pure in-memory.

## When to reach for buntdb

- Local cache for derived data that's cheap to recompute on cold start.
- Test/fixture KV that you don't want to spin up Redis for.
- Sidecar state for a worker (rate limit counters, circuit-breaker state).

For shared state across replicas, **don't use buntdb** — it's per-process. Use `boost-wrapper-cache` with a Redis driver.

## Red flags

| Red flag | Fix |
|---|---|
| `buntdb.Open(path)` with hardcoded path | `buntdb.NewDB(ctx)` + `BOOST_FACTORY_BUNTDB_PATH` |
| Sharing across replicas expecting consistency | Switch to Redis via `boost-wrapper-cache` |
| Forgetting `defer db.Close()` | Add it |
