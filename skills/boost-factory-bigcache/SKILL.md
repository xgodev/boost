---
name: boost-factory-bigcache
description: "Use when constructing an allegro/bigcache in-memory cache via github.com/xgodev/boost/factory/contrib/allegro/bigcache/v3. Covers NewCache + variants and the typical use case (per-process L1 cache in front of Redis or a database, GC-friendly for millions of entries). Triggers on imports under factory/contrib/allegro/bigcache/, on questions about bigcache, L1 caching, or GC pressure from large in-memory caches."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. For typed cache abstraction over multiple backends → `boost-wrapper-cache`.

```go
import bigcachefact "github.com/xgodev/boost/factory/contrib/allegro/bigcache/v3"

c, err := bigcachefact.NewCache(ctx)
if err != nil { log.Fatalf("bigcache: %v", err) }
```

Configure shards, lifetime, max entry size under `boost.factory.bigcache.*` (override `BOOST_FACTORY_BIGCACHE_*`).

## Why bigcache vs sync.Map / a Go map

bigcache stores entries off-heap so the GC doesn't scan them — material when entry count > 100k. For typed access with codecs, prefer `boost-wrapper-cache` with the bigcache driver.

## Red flags

| Red flag | Fix |
|---|---|
| `bigcache.New(...)` directly | `bigcachefact.NewCache(ctx)` |
| Tunables hardcoded | `BOOST_FACTORY_BIGCACHE_*` |
| Used for shared state across replicas | bigcache is per-process; use Redis |
