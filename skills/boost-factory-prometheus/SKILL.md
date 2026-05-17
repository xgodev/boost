---
name: boost-factory-prometheus
description: "Use when wiring Prometheus metric collection in a Go service via github.com/xgodev/boost/factory/contrib/prometheus/client_golang/v1. Covers the typical /metrics endpoint registration pattern (mounted on the main Echo server or on a dedicated multiserver listener). Triggers on imports under factory/contrib/prometheus/client_golang/, on questions about Prometheus metrics in a boost service, or on /metrics endpoint registration."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`. Mount typically via `boost-factory-echo` or `boost-extra-multiserver`.

The factory provides the integration glue between boost and `prometheus/client_golang` — registries, default collectors, and the `promhttp.Handler` wrapper. Concrete API surface is mostly registration helpers; consult the source under `factory/contrib/prometheus/client_golang/v1/` for the symbols available in your boost version.

## Pattern: dedicated /metrics listener via multiserver

```go
import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/xgodev/boost/extra/multiserver"
)

metricsMux := http.NewServeMux()
metricsMux.Handle("/metrics", promhttp.Handler())
metricsSrv := &http.Server{Addr: ":9090", Handler: metricsMux}

ms := multiserver.New(
    multiserver.WithServer("api",     apiSrv),
    multiserver.WithServer("metrics", metricsSrv),
)
ms.Run(ctx)
```

Splitting metrics from the API port means scraping doesn't compete with user traffic and metrics stay accessible if the API saturates.

## Red flags

| Red flag | Fix |
|---|---|
| Metrics on the same port as the public API exposed to the internet | Move to a dedicated listener (private port) via multiserver |
| Custom registry built per request | One registry at startup; all metrics register against it |
| Metrics with high-cardinality labels (per-user-id, per-trace-id) | Cardinality is a Prometheus killer — keep label values bounded |
