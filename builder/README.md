# builder

`builder` generates the doors of an application on top of
[`inject`](../inject): a handler declares with an annotation which door serves
it (a REST route on echo, a gRPC service or the event function), and the
generator writes one fx module per door, a `Modules()` per app and the `main`
of every app. `@Middleware` on a constructor wraps every method of the
interface it provides in `extra/middleware` plugins. One command runs inject
and builder in one pass; the application writes only handlers, plain services
and annotations:

```bash
go generate ./...
go run ./cmd/api
```

## Annotations

### `@REST`: a route on echo

```go
package createorder

// New provides the handler.
// @Provide(index=0)
func New(svc *order.Service) *Handler { return &Handler{svc: svc} }

// Create creates an order.
// @REST(method=POST, path=/orders)
func (h *Handler) Create(c echo.Context) error { ... }
```

```go
package probe

// Live answers the liveness probe.
// @REST(method=GET, path=/live, server=admin)
// @REST(method=HEAD, path=/live, server=admin)
func Live(c echo.Context) error { ... }
```

| Key | Required | Meaning |
|---|---|---|
| `method` | yes | `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD` or `OPTIONS` |
| `path` | yes | echo path (`/orders/:id`) |
| `server` | no | server name, default `rest`; each server is its own echo instance |
| `app` | no | application (binary) that hosts the route, default `api` |

- The handler is `func(echo.Context) error`, a function or a method.
- On a method, the receiver (`*T` or `T`) must be provided by a default
  `@Provide`: the door gets it from the graph.
- `@REST` is repeatable: each line is one route of the same handler.

### `@GRPC`: a service on the gRPC server

```go
package healthcheck

import healthpb "google.golang.org/grpc/health/grpc_health_v1"

// NewServer provides the health service.
// @Provide(index=0)
// @GRPC(register=healthpb.RegisterHealthServer)
func NewServer() *health.Server { ... }
```

| Key | Required | Meaning |
|---|---|---|
| `register` | yes | the `RegisterXServer` function, written as the file writes it (`pkg.RegisterXServer` with `pkg` imported by the file, or a function of the same package) |
| `server` | no | server name, default `grpc` |
| `app` | no | application, default `api` |

`@GRPC` goes on a constructor with exactly one `@Provide(index=N)` without
`name` or `group`; the door calls `register(srv.ServiceRegistrar(), value)`.

### `@Function`: the event handler

```go
package ordercreated

// Handle records the order of the event.
// @Function(adapter=cloudevents, app=worker)
func (h *Handler) Handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) { ... }
```

| Key | Required | Meaning |
|---|---|---|
| `adapter` | yes | `cloudevents`, `kafka`, `nats` or `pubsub` |
| `app` | no | application, default `api` |

- The handler is a [`function.Handler[T]`](../bootstrap/function):
  `func(context.Context, cloudevents.Event) (T, error)`; `T` is its first
  result. A method needs a provided receiver, as in `@REST`.
- An app runs at most one `@Function`. The door includes
  `functionfx.Module[T]` and the adapter module under
  `fx/modules/bootstrap/function/adapter/contrib` (`pubsub` includes v1 and v2,
  each enabled by its own config).

### `@Middleware`: plugins around every method of a service

```go
package stock

// New provides the stock service.
// @Provide(index=0)
// @Middleware(log, prometheus)
func New() Service { return &memory{} }
```

- Plugin names, at least one, no duplicates, the first is the outermost:

| Name | Plugin |
|---|---|
| `log` | [`plugins/local/wrapper/log`](../factory/local/extra/middleware/plugins/local/wrapper/log) |
| `hystrix` | [`plugins/contrib/afex/hystrix-go/v0`](../factory/local/extra/middleware/plugins/contrib/afex/hystrix-go/v0) |
| `datadog` | [`plugins/contrib/datadog/dd-trace-go/v1`](../factory/local/extra/middleware/plugins/contrib/datadog/dd-trace-go/v1) |
| `prometheus` | [`plugins/contrib/prometheus/client_golang/v1`](../factory/local/extra/middleware/plugins/contrib/prometheus/client_golang/v1) |

- It goes on a function with one `@Provide` (with or without `name`, not
  `group`) whose provided result is an exported, non-generic interface. Every
  method (embedded ones included) takes `context.Context` first and returns
  `error`, `(R, error)` or `R`, and uses only exported types.
- builder generates a decorator of the interface and applies it with a root
  `fx.Decorate`, so every consumer in the app receives it, whatever module it
  lives in. Each method runs `next` in the wrapper that fits its results
  (`ErrorWrapper`, `AnyErrorWrapper[R]`, `AnyWrapper[R]`); the service's
  methods stay plain.
- The wrapper of a method is named `<package>.<interface>.<method>`,
  lowercased (`stock.service.available`): it is the hystrix command, the span
  and metric name and the config path of each plugin. Config only turns one
  off: `boost.factory.local.extra.middleware.plugins.<plugin>.<name>.enabled`
  (default `true`).
- An app whose doors never build the constructor gets neither the decorator
  nor its plugins.

App and server names match `^[a-z][a-z0-9]*$` (they become a package, a file,
a function and a config path); an app is not a Go keyword, `main` or
`middleware`, and a server is not `function`.

## Generating

```go
//go:generate go run github.com/xgodev/boost/builder/cmd/builder --boost.builder.output=gen
```

| Key | Env | Default | Description |
|---|---|---|---|
| `boost.builder.path` | `BOOST_BUILDER_PATH` | `.` | directory scanned for the builder and inject annotations |
| `boost.builder.output` | `BOOST_BUILDER_OUTPUT` | `gen` | output directory; must be inside the Go module |
| `boost.builder.cmd` | `BOOST_BUILDER_CMD` | `cmd` | directory that receives `<app>/main.go` for every app; must be inside the Go module |

```
cmd/<app>/main.go                    # boost.Start() and fxfact.NewApp(<app>.Modules()...).Run()
gen/
  module.go                          # inject root module
  github.com/acme/app/.../*_module.go # inject modules, one per @Provide/@Invoke
  builder/<app>/modules.go           # func Modules() []fx.Option
  builder/<app>/<server>_module.go   # one REST or gRPC server and its routes/services
  builder/<app>/function_module.go   # the @Function handler and its adapter
  builder/middleware/<import path>/<func>_middleware.go # Decorate<Func>(): the @Middleware decorator
```

`Modules()` returns the doors sorted by name, then the `@Middleware`
decorators the doors reach and, when the app has servers,
`multiserverfx.Module()`, which starts all of them. Each door includes the
inject modules of its handlers, so only what the app serves is built.
Generated files start with `// Code generated by builder; DO NOT EDIT.`; every
run removes the ones no longer produced, including the main of a removed app.

builder owns the main of every app: a `cmd/<app>/main.go` without that header
fails the generation instead of being overwritten. Other files under `cmd/`
(tests next to the main, other commands) are kept.

### Servers

- `rest` and `grpc` use the shipped
  [echo](../fx/modules/factory/contrib/labstack/echo/v4) and
  [gRPC](../fx/modules/factory/contrib/google.golang.org/grpc/v1/server) fx
  modules, configured under `boost.factory.echo` and
  `boost.factory.grpc.server`.
- Any other name is a separate server configured under
  `<app>.server.<name>` (`api.server.admin.port`): give it its own port.
- A server name serves REST or gRPC, not both.
- Plugins are optional and come from `@Provide`: a `[]boostecho.Plugin`
  (REST) or `[]server.Plugin` (gRPC,
  [`grpc/v1/server`](../factory/contrib/google.golang.org/grpc/v1/server))
  provided without `name` reaches `rest` or `grpc`, and one provided with
  `name=<server>` reaches that server. The door includes the provider's inject
  module, so the generated main needs nothing else:

```go
// Plugins tags the admin responses.
// @Provide(name=admin, index=0)
func Plugins() []boostecho.Plugin { ... }
```

## Validation

The command fails, without touching the output directory, when:

- `method` is missing or unknown, or `path` is missing;
- the same `method`+`path` is declared twice in one server, or a gRPC
  `register` twice;
- the handler signature does not match the door;
- a method handler has no `@Provide` of its receiver;
- `@GRPC` has no single default `@Provide`, or `register` does not resolve
  to an exported function;
- an app has two `@Function`, or `adapter` is missing or unknown;
- a server is used by REST and gRPC;
- an app or server name is invalid, or an annotation has an unknown key;
- the handler is unexported or in package `main`;
- `@Middleware` lists no plugin, an unknown one or one twice, is repeated,
  is not on a single `@Provide`, decorates a group, a type that is not an
  exported non-generic interface with methods, or a method that breaks the
  rules above, or two constructors decorate the same interface;
- `cmd/<app>/main.go` exists and was not generated by builder;
- any inject validation fails (see [`inject`](../inject#validation)).

## Library

```go
modulePath, moduleDir, _ := inject.FindModule(".")
entries, _ := builder.CollectEntries(".")
g, _ := inject.NewGraphFromEntries(ctx, entries)
spec, _ := builder.NewSpec(entries, g, builder.LoadPackages(".")) // loads the interfaces @Middleware decorates
opts := builder.Options{
	GeneratorOptions: inject.GeneratorOptions{ModulePath: modulePath, ModuleDir: moduleDir, OutputDir: outputDir},
	CmdDir:           cmdDir, // empty: no main
}
doors := builder.NewGenerator(spec, opts)
files, _ := doors.Render(ctx)    // builder/<app>/... and builder/middleware/... relative to OutputDir
mains, _ := doors.RenderCmd(ctx) // <app>/main.go relative to CmdDir
// or doors.Generate(ctx) to write both
```

`Render` and `RenderCmd` do not write anything, which lets a test check that
`gen/` and the mains are up to date.

## Example

[`examples/simple`](./examples/simple) has two apps:

- `api`: REST `POST /orders` and `GET /orders/:id` on `rest`, `GET`/`HEAD
  /live` on `admin` (built with the `X-Served-By: admin` plugin of
  `probe.Plugins`), and the gRPC health service on `grpc`; `POST /orders`
  checks `stock.Service`, whose constructor carries `@Middleware(log,
  prometheus)`;
- `worker`: the `ordercreated` function on the cloudevents adapter.

`cmd/api/main.go` and `cmd/worker/main.go` are generated. Its tests call the
routes and the gRPC service of the generated app, check that the handler
receives the stock decorator and that prometheus counted its call, run the
function handler, and fail when `gen/` or a main is stale.

## Non-goals

- A YAML spec: annotations are the only source, and server settings stay in
  boost config.
- Generating OpenAPI or proto from annotations: the spec in `api/` is the
  source, and code may come from it, never the reverse.
- Transport-agnostic handlers: a handler is native to its door.
- Per-route middleware and GraphQL.
- `cache` and `fallback` in `@Middleware` (they need a cache manager or a
  fallback per method), and different middlewares per method.
