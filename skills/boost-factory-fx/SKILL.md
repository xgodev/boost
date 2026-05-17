---
name: boost-factory-fx
description: "Use when constructing the uber/fx app container in a Go service via github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1. Covers the boost-aware fx.App builder, which wires lifecycle hooks, logger, and signal handling consistent with the rest of boost. Triggers on imports under factory/contrib/go.uber.org/fx/, on questions about fx.New / fx.App construction, or on bootstrapping an fx-driven boost service. For module composition patterns, see boost-fx-modules."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-fx-modules` (composition patterns), `boost-wrapper-log` (the logger fx uses).

The factory wraps `go.uber.org/fx`'s `fx.New` with boost defaults: the app's logger comes from `wrapper/log` so fx events appear in your structured log stream, and the lifecycle is wired into the boot sequence so `OnStart` / `OnStop` hooks run with proper context cancellation.

```go
import (
    "go.uber.org/fx"
    fxfact "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1"
)

func main() {
    boost.Start()

    app := fxfact.New(
        myModule.Module(),
        otherModule.Module(),
        fx.Provide(NewService),
        fx.Invoke(registerRoutes),
    )
    app.Run()
}
```

## When to use this vs raw `fx.New`

Use the factory when you want fx events (provider registration, lifecycle transitions) to flow through `wrapper/log` instead of fx's own stdout writer. For libraries that ship modules but don't run an app, raw `fx.Module(...)` is fine.

## Red flags

| Red flag | Fix |
|---|---|
| `fx.New(...)` directly with `fx.WithLogger(fxevent.NopLogger)` | `fxfact.New(...)` so events flow through `wrapper/log` |
| `app.Run()` before `boost.Start()` | `boost.Start()` is always first |
| Mixing `fxfact.New` with `fx.New` in the same binary | Pick one |
