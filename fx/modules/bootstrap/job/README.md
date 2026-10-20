# fx/modules/bootstrap/job

fx module that runs the provided [`job.Handler`](../../../../bootstrap/job)
once and stops the app, so the process exits when the job is done:

```go
boost.Start()
fxfact.NewApp(
	fx.Provide(func(a *Auditor) job.Handler { return a.Run }),
	jobfx.Module(),
).Run()
```

builder generates this for `@Job`.

## Lifecycle

- **OnStart** runs the handler in a goroutine, with a context the stop
  cancels, so `fx.New` returns.
- When the handler returns, the module calls `Shutdowner.Shutdown` with exit
  code 0 on nil; an error is logged at ERROR (`job failed: <error>`) and
  exits 1.
- **OnStop** (SIGTERM, or another hook stopping the app first) cancels the
  handler's context and waits for it to return, up to the stop timeout (fx's
  15s by default). Its error fails the stop, so the app exits 1, and nil lets
  it exit 0; a handler still running at the timeout fails the stop with the
  deadline. A job cut short never reports success unless it says so.

`Module()` returns an empty option after the first call.
