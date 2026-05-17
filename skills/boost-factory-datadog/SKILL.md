---
name: boost-factory-datadog
description: "Use when wiring Datadog APM tracing in a Go service via github.com/xgodev/boost/factory/contrib/datadog/dd-trace-go/v1. Covers NewLogger and NewOptions for the dd-trace-go bridge. Triggers on imports under factory/contrib/datadog/dd-trace-go/, on questions about Datadog tracing, dd-trace-go, or APM integration in a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-log` (logger interop).

```go
import ddfact "github.com/xgodev/boost/factory/contrib/datadog/dd-trace-go/v1"
import "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

opts, _ := ddfact.NewOptions()
tracer.Start(
    tracer.WithLogger(ddfact.NewLogger()),
    // ... apply opts ...
)
defer tracer.Stop()
```

`NewLogger()` returns a `ddtrace.Logger` adapter that funnels Datadog's internal logs through boost's `wrapper/log`, so trace/log routing stays unified.

Configure agent host/port, service name, env, version, sampling under `boost.factory.datadog.*` (override `BOOST_FACTORY_DATADOG_*`). Standard Datadog env vars (`DD_AGENT_HOST`, `DD_SERVICE`, ...) also work via koanf.

## Datadog vs OpenTelemetry

| Pick | When |
|---|---|
| `boost-factory-datadog` | Stack already in Datadog; full APM features (continuous profiler, trace search, RUM) |
| `boost-factory-otel` | Vendor-neutral pipeline (OTLP collector → Tempo/Jaeger/Honeycomb/...) |

Don't mix — one tracing pipeline per service.

## Red flags

| Red flag | Fix |
|---|---|
| `tracer.Start()` without a custom logger | Pass `tracer.WithLogger(ddfact.NewLogger())` so logs are unified |
| Service name / env via raw `os.Getenv("DD_SERVICE")` | Use `BOOST_FACTORY_DATADOG_SERVICE` (registered via `config.Add`) |
| `tracer.Stop()` missing on shutdown | Add `defer tracer.Stop()` after `tracer.Start` |
