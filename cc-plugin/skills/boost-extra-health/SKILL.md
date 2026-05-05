---
name: boost-extra-health
description: "Use when wiring health-check endpoints (/health, /readyz, /livez) in a Go service that imports github.com/xgodev/boost/extra/health. Covers the Checker interface, aggregator pattern, HTTP endpoint registration, and the strict liveness-vs-readiness distinction (liveness = process up; readiness = downstream healthy). Triggers on imports of extra/health, on questions about health checks, readiness vs liveness, downstream dependency probes, or on a Kubernetes deployment with cascading health failures."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-start` — `boost.Start()` first.
- `boost-factory-echo` — endpoints typically attach to the Echo server.

## Liveness vs readiness — keep them separate

| Endpoint | Answers | Fails when |
|---|---|---|
| `/livez` (or `/health`) | "Is the process alive?" | Process should be killed and restarted |
| `/readyz` | "Should I receive traffic right now?" | A downstream is unhealthy; don't pull this pod from rotation, just stop sending requests |

Cascading failure mode to avoid: a transient downstream blip flips `/health` to 500, Kubernetes kills the pod, the new pod hits the same downstream, gets killed, replicas churn → total outage from a recoverable hiccup. Liveness probes must NOT depend on downstream state.

## Wiring

```go
import (
    "github.com/xgodev/boost/extra/health"
)

// Liveness: cheap, deterministic, no downstream calls.
srv.GET("/livez", func(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
})

// Readiness: checks every registered downstream checker.
checkers := []health.Checker{
    health.NewDBChecker(db),
    health.NewPubSubChecker(pb),
    health.NewRedisChecker(redisClient),
}
srv.GET("/readyz", health.AggregateHandler(checkers...))
```

Each checker implements:

```go
type Checker interface {
    Name() string
    Check(ctx context.Context) error
}
```

A nil error from every checker → 200; any error → 503 with the failing checker's name in the response.

## Implementation tips

- **Bound the timeout** on each checker (`context.WithTimeout`) so a hung downstream doesn't hang the readiness response.
- **Cache the result** for a few seconds if checkers are expensive (DB ping, full Redis round-trip) — Kubernetes typically polls every 10s, you don't need fresh data per request.
- **Don't include the publisher** in liveness — publishing to a downed Pub/Sub topic should not kill the pod.

## Red flags

| Red flag | Fix |
|---|---|
| Liveness probe runs downstream checks (DB, Redis, Pub/Sub) | Move downstream checks to readiness; liveness stays static |
| Single `/health` endpoint serving both probes | Split into `/livez` (static) and `/readyz` (dependency-aware) |
| Checker without a context timeout | Wrap with `context.WithTimeout(ctx, 2*time.Second)` |
| Aggregator returns 200 even when one checker failed | Use `health.AggregateHandler` (any failure → 503) |
| Health endpoint reads `os.Getenv("HEALTHY")` to flip state | Reflect actual state from real checkers; don't gate via env |
