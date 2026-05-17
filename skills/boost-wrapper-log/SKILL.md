---
name: boost-wrapper-log
description: "Use when writing or reviewing logging code in a Go service that imports github.com/xgodev/boost/wrapper/log, or when migrating away from log.New / zap.NewProduction / zerolog.New / slog.New in a boost service. Triggers on imports of github.com/xgodev/boost/wrapper/log, on questions about log.FromContext, WithField, WithTypeOf, structured logging, or when a code review spots a third-party logger constructed in main.go of a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start` — `log.FromContext` returns a no-op until `boost.Start()` runs.

## Always pull the logger from context

```go
import "github.com/xgodev/boost/wrapper/log"

logger := log.FromContext(ctx).
    WithField("component", "orders").
    WithField("order_id", orderID)

logger.Info("processing order")
logger.WithError(err).Warn("validation failed")
```

The request- or call-scoped context carries enrichment that previous middlewares added (request ID, tracing, etc.). Pulling from context preserves that chain; constructing your own logger throws it away.

## API surface

| API | When |
|---|---|
| `log.FromContext(ctx)` | Always — request- or call-scoped |
| `log.WithField(k, v)` / `WithError(err)` | Enrich the logger with structured fields |
| `log.WithTypeOf(*p)` | Inside library/driver code to stamp the type name |
| `log.Fatal/Fatalf` | Process-fatal init failures only — never inside a handler |

The configured backend (zap, zerolog, logrus) is picked at `boost.Start` based on `boost.factory.<backend>.console.level` config. Switching backends doesn't require code changes — the wrapper is the abstraction.

## Red flags

| Red flag | Fix |
|---|---|
| `log.New(...)`, `zap.NewProduction()`, `zerolog.New(os.Stdout)`, `slog.New(...)` | Replace with `log.FromContext(ctx)` |
| Constructing a logger as a struct field, holding across requests | Use `log.FromContext` per call so per-request enrichment flows through |
| Calling `log.Fatal` inside a request handler or worker callback | Return an error instead — fatal kills the whole process |
| Logging without structured fields (`logger.Infof("user %s logged in", id)`) | Use `WithField("user_id", id).Info("user logged in")` so log aggregators index it |
