# fx/modules/bootstrap/function

fx module that runs a `bootstrap/function` handler through the adapter
commands provided in the `boostrap.function.adapters` group
(`BSFunctionAdaptersGroupKey`) — cloudevents HTTP, NATS, Pub/Sub, Kafka, AWS Lambda.

```go
fxfact.NewApp(
	fx.Provide(func() function.Handler[*event.Event] { return handle }),
	cehttp.Module[*event.Event](),
	functionfx.Module[*event.Event](middlewares),
).Run()
```

`T` must be `*event.Event` or `[]*event.Event`.

The function runs the middlewares given to `Module`, then the
`[]middleware.AnyErrorMiddleware[T]` provided under the name
`MiddlewaresName` (`bootstrap.function.middlewares`), when there is one. That
is how a builder app passes the `middlewares` of its `@Function`:

```go
fx.Provide(fx.Annotate(
	func() []middleware.AnyErrorMiddleware[*event.Event] { return mws },
	fx.ResultTags(`name:"`+functionfx.MiddlewaresName+`"`),
))
```

## Lifecycle

- **OnStart** — `Function.Run` runs in a goroutine, so `fx.New` returns and
  `app.Run()` listens for SIGINT/SIGTERM.
- **OnStop** — the context passed to `Run` (and to the adapter command) is
  cancelled; the hook waits for `Run` to return, bounded by the stop context
  (fx's stop timeout), and returns the context error if it does not.
- **Run ends on its own** (no stop requested) — the app is shut down with
  exit code 0 on success, or 1 when `Run` returned an error (logged).
