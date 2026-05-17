---
name: boost-factory-cobra
description: "Use when building a multi-command CLI in a Go service via github.com/xgodev/boost/factory/contrib/spf13/cobra/v1. Covers NewCommand for composing a root command with subcommands, including how the function bootstrap (bootstrap/function/function.go) uses cobra internally to drive its CLI shape. Triggers on imports under factory/contrib/spf13/cobra/, on questions about CLI flag handling or subcommand wiring in a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`. For event-driven function CLI shape (`fn.Run` uses cobra internally) → `boost-bootstrap-function`.

```go
import cobrafact "github.com/xgodev/boost/factory/contrib/spf13/cobra/v1"
import co "github.com/spf13/cobra"

root := &co.Command{Use: "myapp"}
sync := &co.Command{Use: "sync", RunE: runSync}
backfill := &co.Command{Use: "backfill", RunE: runBackfill}

cmd := cobrafact.NewCommand(root, sync, backfill)
cmd.Execute()
```

`NewCommand(root, subcommands...)` wires boost-aware defaults (config flag, version flag, structured logging) onto the root before attaching subcommands.

## When to use vs when to skip

Use cobra when the binary genuinely has multiple modes (sync, backfill, migrate, dev tools). For single-purpose binaries (a function or an HTTP API), skip cobra — `boost.Start()` + a plain `main` is simpler.

## Red flags

| Red flag | Fix |
|---|---|
| Hand-rolling cobra with no boost integration | `cobrafact.NewCommand(root, ...)` |
| Defining subcommands inside `init()` of multiple files | Compose them in `main`, not via package-init side effects |
| Reading flags via `os.Args` directly | Use cobra's `cmd.Flags()` |
