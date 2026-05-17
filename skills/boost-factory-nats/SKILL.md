---
name: boost-factory-nats
description: "Use when constructing a raw NATS connection via github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1. Covers NewConn + variants and the Plugin slot. Use this skill for the FACTORY layer (raw *nats.Conn for direct subscribe/publish/JetStream); use boost-bootstrap-adapter-nats for event-handler subscriber wiring with the middleware chain."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`. For event-handler integration → `boost-bootstrap-adapter-nats`.

```go
import natsfact "github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1"

nc, err := natsfact.NewConn(ctx)
if err != nil { log.Fatalf("nats: %v", err) }
defer nc.Close()
```

Configure servers, credentials, name, reconnect policy under `boost.factory.nats.*` (override `BOOST_FACTORY_NATS_*`).

## Factory vs adapter

| Use case | Reach for |
|---|---|
| Direct subscription / publish / request-reply, JetStream | `boost-factory-nats` (raw conn) |
| Event handler over CloudEvents semantics with middleware chain | `boost-bootstrap-adapter-nats` |

## Red flags

| Red flag | Fix |
|---|---|
| `nats.Connect(url)` directly | `natsfact.NewConn(ctx)` |
| Servers/creds via `os.Getenv` | `BOOST_FACTORY_NATS_*` |
| Forgetting `defer nc.Close()` | Add it |
