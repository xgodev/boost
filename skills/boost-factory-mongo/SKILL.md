---
name: boost-factory-mongo
description: "Use when constructing a MongoDB connection in a Go service via github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1 or v2. Covers the constructor trio (NewConn / NewConnWithOptions / NewConnWithConfigPath), the Plugin slot, and the v1 vs v2 path choice. Triggers on imports under factory/contrib/go.mongodb.org/mongo-driver/, on questions about MongoDB connections, NewConn, or which mongo-driver major to use."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`.

## Canonical examples (ship with boost)

- `factory/contrib/go.mongodb.org/mongo-driver/v1/examples/ping/main.go`
- `factory/contrib/go.mongodb.org/mongo-driver/v2/examples/ping/main.go`

Read those before writing new wiring — they are the framework's authoritative shape.

## Choose v1 or v2

Pick the version that matches the mongo-driver major your service depends on. Both expose the same constructor trio.

## Construction

```go
import mongofact "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2"

conn, err := mongofact.NewConn(ctx)
if err != nil { log.Fatalf("mongo: %v", err) }
defer conn.Close(ctx)
```

Configure via `boost.factory.mongo.*` (override `BOOST_FACTORY_MONGO_*`). Multi-database: `mongofact.ConfigAdd("boost.factory.mongo.<name>")` per logical DB + `NewConnWithConfigPath`.

## Red flags

| Red flag | Fix |
|---|---|
| `mongo.Connect(ctx, ...)` directly from upstream SDK | `mongofact.NewConn(ctx)` |
| URI via `os.Getenv` | `BOOST_FACTORY_MONGO_*` |
| Forgetting `defer conn.Close(ctx)` | Add it |
