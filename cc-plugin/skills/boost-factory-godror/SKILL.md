---
name: boost-factory-godror
description: "Use when connecting to Oracle DB from a Go service via github.com/xgodev/boost/factory/contrib/godror/godror/v0. Covers NewDB + variants returning *sql.DB, godror tunables (DSN, pool, statement cache), the wrapper/sql Plugin slot, and the Oracle Instant Client runtime requirement. Triggers on imports under factory/contrib/godror/, on questions about Oracle DB in a boost service, or on NewDB."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. Layout/trio convention → `boost-maintainer`.

```go
import oracle "github.com/xgodev/boost/factory/contrib/godror/godror/v0"

db, err := oracle.NewDB(ctx)
if err != nil { log.Fatalf("oracle: %v", err) }
defer db.Close()
```

Returns `*sql.DB` backed by godror. Configure DSN, user, password, pool sizes under `boost.factory.godror.*` (override `BOOST_FACTORY_GODROR_*`). `plugins ...sqll.Plugin` for cross-cutting concerns.

> **Runtime requirement:** godror needs Oracle Instant Client present at runtime — verify your container/host has the libs before wiring. Build will succeed without them; the connection at startup will fail.

## Red flags

| Red flag | Fix |
|---|---|
| `sql.Open("godror", dsn)` with hand-built DSN | `oracle.NewDB(ctx)` |
| Credentials via `os.Getenv` | `BOOST_FACTORY_GODROR_*` |
| Skipping Oracle Instant Client setup in deployment | Add to base image / install step |
| Forgetting `defer db.Close()` | Add it |
