---
name: boost-factory-cassandra
description: "Use when constructing a Cassandra session via github.com/xgodev/boost/factory/contrib/gocql/gocql/v1. Covers NewSession + variants and the canonical health/ping shapes shipped under the factory's examples/. Triggers on imports under factory/contrib/gocql/, on questions about Cassandra/Scylla in a boost service, or NewSession."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`.

## Canonical examples (ship with boost)

- `factory/contrib/gocql/gocql/v1/examples/ping/main.go`
- `factory/contrib/gocql/gocql/v1/examples/health/main.go`

The `health` example is the reference for wiring this factory into `boost-extra-health` checkers.

## Construction

```go
import gocqlfact "github.com/xgodev/boost/factory/contrib/gocql/gocql/v1"

session, err := gocqlfact.NewSession(ctx)
if err != nil { log.Fatalf("cassandra: %v", err) }
defer session.Close()
```

Configure hosts, keyspace, consistency, timeouts under `boost.factory.gocql.*` (override `BOOST_FACTORY_GOCQL_*`).

## Red flags

| Red flag | Fix |
|---|---|
| `gocql.NewCluster(...).CreateSession()` direct | `gocqlfact.NewSession(ctx)` |
| Hosts/keyspace via `os.Getenv` | `BOOST_FACTORY_GOCQL_*` |
| Forgetting `defer session.Close()` | Add it |
| Hand-rolled health check instead of mirroring `examples/health` | Mirror the example's shape so checker conventions stay uniform |
