---
name: boost-factory-freecache
description: "Use when constructing a coocood/freecache in-memory cache via github.com/xgodev/boost/factory/contrib/coocood/freecache/v1. Covers NewCache + variants and the typical use case (size-bounded L1 cache with LRU eviction). Triggers on imports under factory/contrib/coocood/freecache/, on questions about freecache or strict-memory-cap in-process caches."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. For typed cache abstraction → `boost-wrapper-cache`.

```go
import freecachefact "github.com/xgodev/boost/factory/contrib/coocood/freecache/v1"

c, err := freecachefact.NewCache(ctx)
if err != nil { log.Fatalf("freecache: %v", err) }
```

Configure size (in bytes) under `boost.factory.freecache.size` (override `BOOST_FACTORY_FREECACHE_SIZE`). Once full, freecache evicts oldest entries — predictable memory ceiling.

## bigcache vs freecache

| Pick | When |
|---|---|
| bigcache | Tunable shard count, large entries; lifetime-based eviction |
| freecache | Strict memory cap; LRU eviction; simpler config |

Both are GC-friendly. Pick whichever a sibling service in your monorepo already uses.

## Red flags

| Red flag | Fix |
|---|---|
| `freecache.NewCache(size)` directly | `freecachefact.NewCache(ctx)` |
| Mixing freecache and bigcache in the same binary | Pick one |
| Used for shared state across replicas | Per-process; use Redis |
