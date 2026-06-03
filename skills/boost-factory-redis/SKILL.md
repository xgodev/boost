---
name: boost-factory-redis
description: "Use when constructing a Redis client (single or cluster) via github.com/xgodev/boost/factory/contrib/redis/go-redis/v9. Covers NewClient and NewClusterClient + the canonical shapes shipped under examples/{client,cluster}/. Use this skill for the FACTORY layer (raw *redis.Client / *redis.ClusterClient); use boost-wrapper-cache for the higher-level Manager[T] cache abstraction over a Redis driver."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. For cache-abstraction usage on top of Redis → `boost-wrapper-cache`.

## Canonical examples (ship with boost)

- `factory/contrib/redis/go-redis/v9/examples/client/main.go` — single instance
- `factory/contrib/redis/go-redis/v9/examples/cluster/main.go` — cluster mode

Pick whichever matches your deployment topology.

## Single instance

```go
import redisfact "github.com/xgodev/boost/factory/contrib/redis/go-redis/v9"

rdb, err := redisfact.NewClient(ctx)
if err != nil { log.Fatalf("redis: %v", err) }
defer rdb.Close()
```

## Cluster

```go
rdb, err := redisfact.NewClusterClient(ctx)
```

Configure under `boost.factory.redis.*` (override `BOOST_FACTORY_REDIS_*`).

## Config keys (single vs cluster differ — don't mix them up)

| Topology | Endpoint key | Env override |
|---|---|---|
| Single (`NewClient`) | `boost.factory.redis.client.addr` (one `host:port`) | `BOOST_FACTORY_REDIS_CLIENT_ADDR` |
| Cluster (`NewClusterClient`) | `boost.factory.redis.cluster.addrs` (comma-separated `host:port` list) | `BOOST_FACTORY_REDIS_CLUSTER_ADDRS` |

Common knobs (same prefix): `...dialTimeout`, `...readTimeout`,
`...writeTimeout`, `...maxRetries`, `...minRetryBackoff`,
`...maxRetryBackoff` → `BOOST_FACTORY_REDIS_DIALTIMEOUT`, etc.

The endpoint must be a **resolvable** address from inside the runtime.
In Kubernetes that's the **Service FQDN**
(`<svc>.<ns>.svc.cluster.local:6379`), never a node/master DNS or a
hardcoded IP — the latter fails intermittently with
`dial tcp: lookup ... no such host`. A single-instance service that sets
`...cluster.addrs` (or vice versa) silently connects to nothing.

## Factory vs cache abstraction

| Use case | Reach for |
|---|---|
| Pub/Sub on Redis, streaming, scripting, low-level commands | `boost-factory-redis` (raw client) |
| Typed cache over a value type with codec + plugins | `boost-wrapper-cache` |

## Red flags

| Red flag | Fix |
|---|---|
| `redis.NewClient(&redis.Options{...})` directly | `redisfact.NewClient(ctx)` |
| Cluster vs single decided at runtime by env-sniff | Pick at deploy time; keep a single example shape per service |
| Forgetting `defer rdb.Close()` | Add it |
