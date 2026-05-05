---
name: boost-factory-echo
description: "Use when writing or reviewing Go HTTP API services that import github.com/xgodev/boost/factory/contrib/labstack/echo/v4. Covers the canonical echo.NewServer construction, the required plugin set (recover, requestid, log, restresponse, error_handler, cors, gzip), the strict plugin order (restresponse before error_handler), and graceful shutdown wrapping. Triggers on imports under factory/contrib/labstack/echo, on questions about echoserver.NewServer, ErrorHandlerJSON vs ErrorHandlerString, REST type, or on creating new HTTP endpoints in a boost service."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:**
- `boost-start` — `boost.Start()` first.
- `boost-wrapper-log` — request logger via context.
- `boost-model-errors` — handler errors are typed and routed by `error_handler`.

## Canonical `main.go`

```go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/xgodev/boost"
    echoserver "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
    "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/extra/error_handler"
    logplugin "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
    restresponse "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/model/restresponse"
    recoverplugin "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/recover"
    "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/requestid"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    boost.Start()

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    srv, err := echoserver.NewServer(ctx,
        recoverplugin.Register, // panics → 500, no process death
        requestid.Register,     // X-Request-ID propagation
        logplugin.Register,     // request access log via boost logger
        restresponse.Register,  // sets Type=REST so error_handler emits JSON
        error_handler.Register, // model/errors → HTTP status mapping
    )
    if err != nil {
        log.FromContext(ctx).WithError(err).Fatal("failed to build echo server")
    }

    srv.GET("/health", healthHandler)
    srv.POST("/orders", orderHandler.Create)

    go srv.Serve(ctx)
    <-ctx.Done()

    shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancelShutdown()
    srv.Shutdown(shutdownCtx)
}
```

## Required plugin set

| Plugin | Purpose | Skip when |
|---|---|---|
| `native/recover` | Panic → 500 | Never |
| `native/requestid` | X-Request-ID | Never |
| `local/wrapper/log` | Access log via boost logger | Never |
| `local/model/restresponse` | Sets `Type=REST` for `error_handler` JSON mode | Never (REST APIs) |
| `extra/error_handler` | Maps `model/errors.*` to HTTP status + JSON envelope | Never |
| `native/cors` | Browser clients | Internal-only services |
| `native/gzip` | Response compression | gRPC / streaming |

## Plugin order constraint

`restresponse.Register` MUST come before `error_handler.Register`. Reason: the error handler picks `ErrorHandlerJSON` only when the server has `Type=REST` set on it. Without `restresponse` first, you get `ErrorHandlerString` (text/plain) — production symptom: 4xx/5xx responses come back as `"some error"` text instead of the JSON envelope clients expect.

## Graceful shutdown

```go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer cancel()

go srv.Serve(ctx)        // blocking; goroutine so we can react to ctx.Done()
<-ctx.Done()

shutdownCtx, _ := context.WithTimeout(context.Background(), 15*time.Second)
srv.Shutdown(shutdownCtx) // bounded drain with FRESH ctx
```

Use a fresh context for `Shutdown` — passing the cancelled parent makes it return immediately.

## Red flags

| Red flag | Fix |
|---|---|
| `echo.New()` directly | `echoserver.NewServer(ctx, ...)` |
| `error_handler.Register` before `restresponse.Register` | Reorder |
| `srv.Serve(ctx)` called inline (not in a goroutine) followed by `<-ctx.Done()` | Wrap `Serve` in a goroutine |
| `echo.NewHTTPError(...)` in a handler | Return `bootsterrors.<Type>(err, "...")` (see `boost-model-errors`) |
| `Shutdown(ctx)` reusing the cancelled parent context | Use a fresh `context.WithTimeout` |
