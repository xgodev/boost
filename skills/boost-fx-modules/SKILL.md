---
name: boost-fx-modules
description: "Use when scaling a Go service beyond a single main.go to module-composed dependency injection via github.com/xgodev/boost/fx/modules (uber-go/fx-based). Covers Module() functions returning fx.Option, the fx.Provide / fx.Invoke wiring patterns, the group:\"...\" tag for collecting middleware contributions, optional:\"true\" for opt-in deps, and when fx is overkill (small services). Triggers on imports under fx/modules/, on questions about uber/fx, fx.New, fx.Provide, fx.Invoke, group annotations, or on refactoring a tangled main.go into composable modules."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`. Each fx module typically wires components covered by a sibling skill (`boost-factory-echo`, `boost-bootstrap-function`, etc.).

## When to reach for fx

| Service shape | Recommendation |
|---|---|
| 1 main.go with handlers + dependencies wired by hand | Don't use fx — manual wiring is clearer |
| Multiple binaries sharing the same DB / pubsub / cache wiring | Extract those into fx modules; reuse across binaries |
| 5+ services in a monorepo with the same boost-startup boilerplate | Use fx modules; standardize the boilerplate |

The Carrefour `br-digicomm-catalog-core` repo follows the latter pattern: shared modules (DB, Pub/Sub client, publisher, logger) live there and every `*-catalog-*` service composes them via `fx.Options(catalogcore.PubsubModule(), catalogcore.DBModule(), ...).Run()`.

## Module shape

```go
package catalogcore

import (
    "go.uber.org/fx"
    "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1"
)

func PubsubModule() fx.Option {
    return fx.Module("pubsub",
        fx.Provide(pubsub.NewClient),
        fx.Invoke(registerLifecycleHooks),
    )
}

func registerLifecycleHooks(lc fx.Lifecycle, c *pubsub.Client) {
    lc.Append(fx.Hook{
        OnStop: func(ctx context.Context) error { return c.Close() },
    })
}
```

`fx.Provide(constructor)` makes the constructor's return value available as a dependency. `fx.Invoke(fn)` runs `fn` at startup and threads in whatever it asks for.

## Composition in `main`

```go
func main() {
    boost.Start()

    app := fx.New(
        catalogcore.PubsubModule(),
        catalogcore.DBModule(),
        catalogcore.PublisherModule(),
        fx.Provide(NewOrderService, NewOrderHandler),
        fx.Invoke(registerHTTPRoutes),
    )
    app.Run()
}
```

`app.Run()` blocks until SIGTERM, then runs every `OnStop` hook in reverse-of-start order.

## Group annotations — for middleware contributions

When multiple modules each contribute a piece to a chain (e.g., function middlewares), use group tags:

```go
type Middlewares struct {
    fx.In
    Items []middleware.AnyErrorMiddleware[*cloudevents.Event] `group:"bootstrap.function.middleware"`
}

func RecoveryModule() fx.Option {
    return fx.Provide(
        fx.Annotate(
            func() middleware.AnyErrorMiddleware[*cloudevents.Event] { return rm.NewRecovery[*cloudevents.Event]() },
            fx.ResultTags(`group:"bootstrap.function.middleware"`),
        ),
    )
}
```

Anyone consuming `Middlewares.Items` gets the full set without knowing which modules contributed.

## Optional deps — `optional:"true"`

When a module's behavior changes based on whether something is provided:

```go
type FunctionDeps struct {
    fx.In
    Publisher *publisher.Publisher[*cloudevents.Event] `optional:"true"`
}
```

If no module provides `*publisher.Publisher[*cloudevents.Event]`, `FunctionDeps.Publisher` is nil and the consuming code degrades gracefully.

## Red flags

| Red flag | Fix |
|---|---|
| Single-binary service with 200-line main.go using fx | Strip fx; manual wiring is clearer |
| `fx.New(...).Run()` called BEFORE `boost.Start()` | `boost.Start()` is always first |
| Constructor returning a concrete type when consumers need an interface | Return the interface so swapping implementations doesn't break the graph |
| `fx.Invoke` doing real business logic | Use `fx.Invoke` only for wiring (route registration, lifecycle hooks); business logic lives in the constructed types |
| Per-service copy-paste of the same fx wiring | Extract into a shared module package (e.g., `<orgname>/<core>/fxmodules`) |
