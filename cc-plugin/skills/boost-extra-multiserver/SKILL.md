---
name: boost-extra-multiserver
description: "Use when a single Go binary needs to expose more than one network listener (e.g., main HTTP API + gRPC + Prometheus /metrics + health side-car) via github.com/xgodev/boost/extra/multiserver. Covers coordinated start/stop, per-server lifecycle, fail-fast on any listener error, and graceful drain of all listeners on SIGTERM. Triggers on imports of extra/multiserver, on questions about multi-listener boost services, or on a binary that needs HTTP + gRPC simultaneously."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-start` — `boost.Start()` first.
- `boost-factory-echo` — the typical primary listener.

## When you need it

A single binary listening on multiple ports — for example:
- Main API on `:8080` (Echo HTTP).
- Prometheus scraper on `:9090` (`/metrics`).
- Internal gRPC on `:9000`.
- Readiness side-car on `:8081` (so the readiness port is independent of the API port).

Without `multiserver`, you end up writing ad-hoc goroutine + WaitGroup orchestration per listener. `multiserver` standardizes start order, fail-fast on any listener error, and SIGTERM-driven coordinated shutdown.

## Wiring

```go
import (
    "github.com/xgodev/boost/extra/multiserver"
)

ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer cancel()

ms := multiserver.New(
    multiserver.WithServer("api",     apiSrv),       // primary HTTP
    multiserver.WithServer("metrics", metricsSrv),   // Prometheus
    multiserver.WithServer("grpc",    grpcSrv),      // gRPC
)

if err := ms.Run(ctx); err != nil {
    log.FromContext(ctx).WithError(err).Fatal("multiserver exited with error")
}
```

`ms.Run(ctx)` blocks until ctx is cancelled OR any registered server returns a fatal error. On exit, every server's `Shutdown` is called with a bounded drain context.

## Failure semantics

- If any listener fails to bind at startup → all servers stop, `Run` returns the error.
- If a listener panics or returns mid-flight → coordinated drain (the rest get a chance to finish in-flight requests before exit).
- SIGTERM on the parent ctx → bounded drain across all listeners in parallel.

## Red flags

| Red flag | Fix |
|---|---|
| Hand-rolled `sync.WaitGroup` + N `go server.Serve` calls | Use `multiserver.New` so failure semantics are uniform |
| Forgetting to Shutdown one of the listeners on signal | Let `multiserver` drive shutdown for all of them |
| Different drain timeouts per listener | Set one drain timeout on the multiserver options; consistency beats micro-tuning here |
| Mixing `multiserver` with manual `srv.Serve` for one of the servers | All listeners under `multiserver` or none — the failure-fast guarantee depends on it |
