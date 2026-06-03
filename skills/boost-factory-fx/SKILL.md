---
name: boost-factory-fx
description: "Use when constructing the uber/fx app container in a Go service via github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1. Covers the boost-aware fx.App builder, which wires lifecycle hooks, logger, and signal handling consistent with the rest of boost. Triggers on imports under factory/contrib/go.uber.org/fx/, on questions about fx.New / fx.App construction, or on bootstrapping an fx-driven boost service. For module composition patterns, see boost-fx-modules."
license: MIT
metadata:
  author: jpfaria
  version: "0.2.0"
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

## Fail loud at boot

A constructor in `fx.Provide` or an `fx.Invoke` that returns an error makes fx fail the app build and `app.Run()` exit non-zero — and the surfaced reason can be **terse**: a config-validation error returned from a provider can read as a bare `exit 1` with little context after the logger banner. So **validate loudly**: in the provider, log the specific failure (which key, which value) *before* returning the error, and prefer `errors.NewInternal(err, "config: <key>=<value>")` over a bare error. A silent `exit 1` once the logger is up almost always means a provider rejected its config — make sure the provider says which one. (dev-rules LAW 4: fail loud, fail specific.)

## Red flags

| Red flag | Fix |
|---|---|
| `fx.New(...)` directly with `fx.WithLogger(fxevent.NopLogger)` | `fxfact.New(...)` so events flow through `wrapper/log` |
| A boot provider returns a bare error on bad config → opaque `exit 1` | Log the specific bad key/value in the provider before returning; wrap with `errors.NewInternal(err, "config: <key>=<value>")`. A swallowed boot error costs hours to diagnose. |
| `app.Run()` before `boost.Start()` | `boost.Start()` is always first |
| Mixing `fxfact.New` with `fx.New` in the same binary | Pick one |
