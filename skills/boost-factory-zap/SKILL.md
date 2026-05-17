---
name: boost-factory-zap
description: "Use when wiring Uber's zap as the logging backend in a Go service that imports github.com/xgodev/boost/factory/contrib/go.uber.org/zap/v1. Covers zap.NewLogger() construction, log.Set() after boost.Start(), and the boost.factory.zap.* config tree (console/file, level, TEXT/JSON formatter). zap is NOT the boost default — zerolog is — so this skill is for services that explicitly opt into zap. Triggers on imports under factory/contrib/go.uber.org/zap, on log.Set(zap.NewLogger()), or on BOOST_FACTORY_ZAP_* env vars."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-start` — `boost.Start()` loads config and sets the default logger first.
- `boost-wrapper-log` — the abstraction you actually log through; this skill only swaps the backend.
- `boost-wrapper-config` — `boost.factory.zap.*` namespacing / env override semantics.

## Construction — swap the backend after Start

```go
import (
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/go.uber.org/zap/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    boost.Start()                 // wires zerolog (the default) from config
    log.Set(zap.NewLogger())      // replace it with zap — config-driven, no zap API here

    log.Infof("hello world!!")
}
```

`zap.NewLogger()` returns a `log.Logger` built entirely from the `boost.factory.zap.*` config tree (it `panic`s if options fail to load). You never touch `zap.NewProduction()` / `zapcore` directly — the factory owns construction so level and formatter stay operator-tunable. Order matters: `log.Set` must come **after** `boost.Start()`, otherwise `config.Load()` hasn't run yet.

## Config tree — `boost.factory.zap` (cited from `factory/contrib/go.uber.org/zap/v1/config.go`)

| Key | Default | What |
|---|---|---|
| `.console.enabled` | `true` | console sink on/off |
| `.console.level` | `INFO` | console level |
| `.console.formatter` | `TEXT` | `TEXT` or `JSON` |
| `.file.enabled` | `false` | file sink on/off |
| `.file.level` | `INFO` | file level |
| `.file.path` | `/tmp` | log directory |
| `.file.name` | `application.log` | log filename |
| `.file.maxsize` | `100` | rotation size (MB) |
| `.file.compress` | `true` | gzip rotated files |
| `.file.maxage` | `28` | retention (days) |
| `.file.formatter` | `TEXT` | `TEXT` or `JSON` |

Operators override via env, e.g. `BOOST_FACTORY_ZAP_CONSOLE_FORMATTER=JSON`, `BOOST_FACTORY_ZAP_CONSOLE_LEVEL=DEBUG` (see `examples/boot/main.go`).

## Red flags

| Red flag | Fix |
|---|---|
| `zap.NewProduction()` / `zap.NewDevelopment()` / `zapcore` in app code | `log.Set(zap.NewLogger())` once in `main`, log via `boost-wrapper-log` |
| `log.Set(zap.NewLogger())` before `boost.Start()` | Move it after — config isn't loaded until `Start()` |
| Calling `zap.NewLogger()` per request / storing `*zap.Logger` | Set the backend once; pull request-scoped loggers via `log.FromContext(ctx)` |
| Setting JSON output by constructing a custom encoder | Set `boost.factory.zap.console.formatter=JSON` (env `BOOST_FACTORY_ZAP_CONSOLE_FORMATTER`) |
| Expecting zap without calling `log.Set` | zerolog is the boost default; zap only takes effect after an explicit `log.Set(zap.NewLogger())` |
