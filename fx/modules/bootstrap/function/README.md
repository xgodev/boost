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

## Which adapter runs

`Function.Run` runs the adapter command the first CLI argument names
(`app nats`). With no command named, it runs the one
`boost.bootstrap.function.default` names, or the only adapter in the group:
an app that wires one adapter, as builder does, runs with no argument and no
key. No adapter, or more than one and none chosen, fails the boot with
`NotValid` (exit code 1) naming the adapters and the key, instead of
printing the usage.

## Adapter clients

The `nats`, `cloudevents_nats` and Pub/Sub v2 adapter modules
(`adapter/contrib/nats-io/nats.go/v1`,
`adapter/contrib/cloudevents/sdk-go/v2/contrib/nats-io/nats.go/v1`,
`adapter/contrib/cloud.google.com/pubsub/v2`) come in two forms:

| Function | Client |
|---|---|
| `Module[T]()` | brings its own: the factory fx module of `*nats.Conn` or `*pubsub.Client` |
| `ModuleWithoutClient[T]()` | takes the one the app provides |
| `ModuleWithNamedClient[T](name)` | takes the one the app provides as `name:"<name>"` (nats, cenats, pubsub v2) |

Use `ModuleWithoutClient` when the app already provides that client, such as
the default client a publisher driver or an `@Inject` gets: fx accepts one
provider per type, so `Module` there fails the boot with "already provided".
builder generates it that way, and `ModuleWithNamedClient` for
`@Function(client=x)`. All share one `sync.Once`: the first call wins.

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

The `nats`, `cloudevents_nats` and `gcp_pubsub` (v1 and v2) adapters run on
that context (`Helper.Run(cmd.Context())`) and block until it is done: when
it is cancelled they take no new message, let the handler finish the one it
holds, with a context the stop does not cancel, and return, so SIGTERM stops
the app with exit code 0 well before the stop timeout. NATS drains its
subscriptions; Pub/Sub gives back a failed message instead of retrying it. A
NATS subject it cannot subscribe to, or a Pub/Sub subscription it cannot
receive on (`pubsub: could not receive on <name>`, e.g. one that does not
exist), fails `Run`, so the boot fails with exit code 1. `Helper.Start()` is
deprecated: it runs with a context that is never cancelled.

The Pub/Sub v1 adapter module provides the v1 `*pubsub.Client`
(`fx/modules/factory/contrib/cloud.google.com/pubsub/v1.ClientModule()`)
whatever `boost.factory.gcp.pubsub.version` says, so
`boost.bootstrap.function.adapter.pubsub.version=v1` alone turns it on.
