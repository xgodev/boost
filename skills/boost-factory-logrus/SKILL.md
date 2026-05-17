---
name: boost-factory-logrus
description: "Use when wiring Sirupsen logrus as the logging backend in a Go service that imports github.com/xgodev/boost/factory/contrib/sirupsen/logrus/v1. Covers logrus.NewLogger(hooks...) construction (it accepts logrus.Hook values), log.Set() after boost.Start(), the boost.factory.logrus.* config tree, and the TEXT/JSON/CLOUDWATCH formatter selected via .formatterType. logrus is NOT the boost default — zerolog is. Triggers on imports under factory/contrib/sirupsen/logrus, on log.Set(logrus.NewLogger(...)), on logrus hooks, or on BOOST_FACTORY_LOGRUS_* env vars."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-start` — `boost.Start()` loads config and sets the default logger first.
- `boost-wrapper-log` — the abstraction you log through; this skill only swaps the backend.
- `boost-wrapper-config` — `boost.factory.logrus.*` namespacing / env override semantics.

## Construction — swap the backend after Start, optionally with hooks

```go
import (
    "github.com/xgodev/boost"
    lg "github.com/sirupsen/logrus"
    "github.com/xgodev/boost/factory/contrib/sirupsen/logrus/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    boost.Start()                          // wires zerolog (the default)
    log.Set(logrus.NewLogger())            // replace with logrus
    // or, attaching native logrus hooks:
    // log.Set(logrus.NewLogger(sentryHook, datadogHook))

    log.Infof("hello world!!")
}
```

`logrus.NewLogger(hooks ...lg.Hook) log.Logger` is the one factory entrypoint that takes arguments — pass native `logrus.Hook` values for sinks the config tree doesn't model (Sentry, etc.). The formatter is chosen at construction time from `FormatterType()` (`boost.factory.logrus.formatterType`): `CLOUDWATCH` → `formatter/cloudwatch`, `JSON` → `formatter/json`, anything else → `formatter/text`. `log.Set` must come **after** `boost.Start()`.

## Config tree — `boost.factory.logrus` (cited from `factory/contrib/sirupsen/logrus/v1/config.go`)

| Key | Default | What |
|---|---|---|
| `.console.enabled` | `true` | console sink on/off |
| `.console.level` | `INFO` | console level |
| `.file.enabled` | `false` | file sink on/off |
| `.file.level` | `INFO` | file level |
| `.file.path` | `/tmp` | log directory |
| `.file.name` | `application.log` | log filename |
| `.file.maxsize` | `100` | rotation size (MB) |
| `.file.compress` | `true` | gzip rotated files |
| `.file.maxage` | `28` | retention (days) |
| `.time.format` | `2006/01/02 15:04:05.000` | timestamp layout |
| `.formatterType` | `TEXT` | `TEXT`, `JSON`, or `CLOUDWATCH` (picks the formatter package) |

Operators override via env, e.g. `BOOST_FACTORY_LOGRUS_FORMATTERTYPE=JSON`, `BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL=DEBUG`.

## Red flags

| Red flag | Fix |
|---|---|
| `logrus.New()` / `logrus.SetFormatter(...)` in app code | `log.Set(logrus.NewLogger())` once in `main`, log via `boost-wrapper-log` |
| `log.Set(logrus.NewLogger())` before `boost.Start()` | Move it after — config isn't loaded until `Start()` |
| Registering hooks via `logrus.AddHook` on a global logger | Pass them to `logrus.NewLogger(hook1, hook2)` so the boost-built logger owns them |
| Hand-building a JSON/CloudWatch formatter | Set `boost.factory.logrus.formatterType` (`JSON` / `CLOUDWATCH`) — the factory wires the matching `formatter/*` package |
| Expecting logrus without calling `log.Set` | zerolog is the boost default; logrus only takes effect after an explicit `log.Set(logrus.NewLogger())` |
| `logrus.NewLogger()` per request / caching the entry | Set the backend once; pull request-scoped loggers via `log.FromContext(ctx)` |
