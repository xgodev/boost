---
name: boost-wrapper-log-backends
description: "Use when configuring or switching the underlying logger backend in a boost service — picking between zap (default), zerolog, or logrus via github.com/xgodev/boost/wrapper/log/contrib/{go.uber.org/zap,rs/zerolog,sirupsen/logrus}/v1. Covers backend selection at boot via config, the per-backend tunables (format, level, output), and why the choice is invisible to handler code (the wrapper API stays the same). Triggers on imports under wrapper/log/contrib/, on questions about switching from zap to zerolog, on backend-specific config keys (BOOST_FACTORY_ZAP_*, BOOST_FACTORY_ZEROLOG_*), or on log format / level mismatches."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-wrapper-log` — the wrapper API your code calls (`log.FromContext`).
- `boost-wrapper-config` — backend selection is config-driven.

## Backend selection at boot

```go
import (
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/go.uber.org/zap/v1"   // or rs/zerolog/v1, or sirupsen/logrus/v1
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    boost.Start()
    log.Set(zap.NewLogger())   // installs zap as the global backend
    // ... rest of main ...
}
```

After `log.Set`, every `log.FromContext(ctx)` call elsewhere in the service returns a logger backed by zap. Swap to zerolog:

```go
import zerolog "github.com/xgodev/boost/factory/contrib/rs/zerolog/v1"

log.Set(zerolog.NewLogger())
```

Handler code does not change — it still calls `log.FromContext(ctx).WithField(...)` etc.

## Available backends

| Backend | Path | Best for |
|---|---|---|
| zap | `factory/contrib/go.uber.org/zap/v1` | Default; high throughput, structured |
| zerolog | `factory/contrib/rs/zerolog/v1` | Smallest allocations; JSON-first |
| logrus | `factory/contrib/sirupsen/logrus/v1` | Legacy compat with logrus-using libraries |

## Per-backend config

Each backend has its own config root. Examples (zap):

| Key | What |
|---|---|
| `boost.factory.zap.console.level` | `DEBUG` / `INFO` / `WARN` / `ERROR` |
| `boost.factory.zap.console.formatter` | `JSON` / `CONSOLE` (text) |

Override at deploy via `BOOST_FACTORY_ZAP_CONSOLE_LEVEL`, etc. Mirror keys exist for zerolog (`BOOST_FACTORY_ZEROLOG_*`) and logrus (`BOOST_FACTORY_LOGRUS_*`).

Pick **JSON** in production (log aggregators index it) and **CONSOLE** locally (humans read it).

## Red flags

| Red flag | Fix |
|---|---|
| `log.Set(...)` called more than once in the same process | Call once, in `main`, right after `boost.Start()` |
| Calling `log.Set` inside an `init()` — before `boost.Start` configured anything | Move to `main` after `boost.Start()` |
| Importing two backend factories simultaneously and switching at runtime | Pick one per binary; switching is a config / deploy decision, not a runtime one |
| Production binary running with `formatter=CONSOLE` | Set `BOOST_FACTORY_<BACKEND>_CONSOLE_FORMATTER=JSON` in the deployment |
| Mixing the wrapper (`log.FromContext`) and direct backend calls (`zap.L().Info(...)`) | Use the wrapper exclusively — direct calls bypass per-context enrichment |
