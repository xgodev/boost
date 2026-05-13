---
name: boost-factory-pgx
description: "Use when constructing a PostgreSQL *sql.DB via github.com/xgodev/boost/factory/contrib/jackc/pgx/v5. Covers the constructor trio (NewDB / NewDBWithOptions / NewDBWithConfigPath), the wrapper/sql Plugin slot for tracing/metrics, and lifecycle. Triggers on imports under factory/contrib/jackc/pgx/, on questions about pgx, PostgreSQL connection in boost, or NewDB."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. Layout/trio convention → `boost-maintainer`.

```go
import pgx "github.com/xgodev/boost/factory/contrib/jackc/pgx/v5"

db, err := pgx.NewDB(ctx)
if err != nil { log.Fatalf("pgx: %v", err) }
defer db.Close()
```

Returns `*sql.DB` (database/sql interface) backed by jackc/pgx v5. Configure host, port, database, user, password, sslMode, pool sizes under `boost.factory.pgx.*` (override `BOOST_FACTORY_PGX_*`).

Multi-DB pattern:

```go
pgx.ConfigAdd("boost.factory.pgx.orders")
pgx.ConfigAdd("boost.factory.pgx.analytics")

ordersDB, _    := pgx.NewDBWithConfigPath(ctx, "boost.factory.pgx.orders")
analyticsDB, _ := pgx.NewDBWithConfigPath(ctx, "boost.factory.pgx.analytics")
```

The `plugins ...sqll.Plugin` slot accepts wrappers from `wrapper/sql` (tracing, slow-query log, metrics).

## Red flags

| Red flag | Fix |
|---|---|
| `sql.Open("pgx", dsn)` with hand-built DSN | `pgx.NewDB(ctx)` |
| Connection URL via `os.Getenv` | `BOOST_FACTORY_PGX_*` |
| Pool sizes hardcoded | Tune via `boost.factory.pgx.pool.*` |
| Forgetting `defer db.Close()` | Add it |
