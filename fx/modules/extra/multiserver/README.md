# fx/modules/extra/multiserver

fx module that serves every `multiserver.Server` provided in the
`extra.multiserver.servers` group (`ServersGroupKey`) — e.g. the echo and
gRPC server modules.

```go
fxfact.NewApp(
	echofx.Module(),
	multiserver.Module(),
).Run()
```

## Lifecycle

- **OnStart** — `multiserver.Serve` runs in a goroutine, so `fx.New` returns
  and `app.Run()` listens for SIGINT/SIGTERM.
- **OnStop** — `Shutdown(stopCtx)` is called on every server concurrently,
  then the hook waits for them to drain and for `Serve` to return, bounded
  by the stop context (fx's stop timeout); it returns the context error if
  the drain does not finish in time. The context passed to `Serve` is
  cancelled afterwards.
- **Every server stopped on its own** (no stop requested) — the app is shut
  down with exit code 1. One server down while others run is reported by
  `multiserver.Check` (health), not by shutting down.
- **No servers in the group** — `fx.New` fails with
  `multiserver: no servers configured`.
