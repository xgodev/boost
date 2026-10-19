# builder

`builder` generates the doors of an application on top of
[`inject`](../inject): a handler declares with an annotation which door serves
it (a REST route on echo, a gRPC service, an MCP tool or the event
function), and the
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

### `@MCP`: a tool on the MCP server

```go
package stocktool

// New provides the handler.
// @Provide(index=0)
func New(stock stock.Service) *Handler { return &Handler{stock: stock} }

// CheckStock tells whether an item can be ordered.
// @MCP
func (h *Handler) CheckStock(ctx context.Context, req *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, Output, error) { ... }
```

| Key | Required | Meaning |
|---|---|---|
| `name` | no | tool name (letters, digits, `_`, `-`, `.`, up to 128), default the handler name in snake case (`check_stock`) |
| `server` | no | server name, default `mcp` |
| `app` | no | application, default `api` |

- The handler is a go-sdk
  [`mcp.ToolHandlerFor[In, Out]`](https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk/mcp#ToolHandlerFor):
  `func(context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, Out, error)`.
  The door calls `mcp.AddTool`, which infers the input and output schemas
  from `In` and `Out` and fills the result from `Out`, so the handler may
  return a nil `*mcp.CallToolResult`. A method needs a provided receiver, as
  in `@REST`.
- The tool description is the doc comment of the handler, without its
  leading name: `Tells whether an item can be ordered.`
- The handler lives in the app layer and calls the service; the service
  never imports `mcp`.

### `@Function`: the event handler

```go
package ordercreated

// Handle records the order of the event.
// @Function(adapter=[cloudevents, nats], app=worker, middlewares=[recovery, logger])
func (h *Handler) Handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) { ... }
```

| Key | Required | Meaning |
|---|---|---|
| `adapter` | yes | `azure` (Azure Functions custom handler), `cloudevents`, `cenats` (cloudevents over NATS), `kafka`, `kafkago` (kafka on segmentio/kafka-go), `lambda` (AWS Lambda), `nats` or `pubsub`; or a list `[a, b]` of them, chosen by config |
| `app` | no | application, default `api` |
| `middlewares` | no | `[a, b]`, the outermost first: `ignore_errors`, `logger`, `opentelemetry`, `prometheus`, `publisher`, `recovery` |

- The handler is a [`function.Handler[T]`](../bootstrap/function):
  `func(context.Context, cloudevents.Event) (T, error)`; `T` is its first
  result. A method needs a provided receiver, as in `@REST`.
- An app runs at most one `@Function`. The door includes
  `functionfx.Module[T]` and the adapter module under
  `fx/modules/bootstrap/function/adapter/contrib` (`pubsub` includes v1 and v2;
  `boost.bootstrap.function.adapter.pubsub.version`, `v2` by default, turns
  on one).
- The `nats`, `cenats` and `pubsub` (v2) adapters read with the app's
  default client: the `*nats.Conn` or `*pubsub.Client` of the app's unnamed
  `@Provide`, or else the one inject generates (`DefaultConnModule()`,
  `DefaultClientModule()` under `gen/client/...`), the same one the publisher
  and an `@Inject` without `name` get. The door wires the adapter's
  `ModuleWithoutClient`, so a function that reads from a broker and
  publishes to it boots with one client.
- `client=orders` reads with the named client `orders` instead: the app's
  `@Provide(name=orders)` of that type, or else the one inject generates
  (`OrdersConnModule()`, keys `boost.factory.nats.clients.orders.*`), the
  same one `@Inject(name=orders)` and the publisher `orders` get. The door
  wires the adapter's `ModuleWithNamedClient("orders")`. The name follows
  the client names (`[a-z][a-zA-Z0-9]*`), and the adapters must include one
  that reads with a client (`nats`, `cenats`, `pubsub`), else `NotValid`.
  The pubsub v1 adapter keeps its own client.
- With a list, `adapter=[pubsub, nats]`, config chooses the adapter:
  `<app>.function.adapter` (env `<APP>_FUNCTION_ADAPTER`), default the first
  name. The door includes only the chosen adapter module, so the others and
  their clients are never built; an unknown name fails the boot with
  `NotValid` naming the key, the value and the valid names. The adapters'
  own keys (`.enabled`, the pubsub `version`) stay as they are. The app runs
  the chosen adapter with no CLI argument and no
  `boost.bootstrap.function.default`.
- The pubsub adapter on v1 (`boost.bootstrap.function.adapter.pubsub.version=v1`)
  reads with its own v1 `*pubsub.Client`, built from the
  `boost.factory.gcp.pubsub` keys; the app's default client is v2.
- Each middleware is the fx module under
  `fx/modules/bootstrap/function/middleware` of the
  [`bootstrap/function/middleware`](../bootstrap/function/middleware) of that
  name; the door hands them to `functionfx.Module` under
  `functionfx.MiddlewaresName`. One turned off in config
  (`boost.bootstrap.function.middleware.<name>.enabled`) is skipped.
  `publisher` needs the app's unnamed `publisher.Driver`: the door builds the
  `@Provide` constructor that returns it, and generation fails with NotFound
  when there is none (fx would skip the middleware silently).

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
- On an implementation of a set (see
  [`inject` -- Default implementation](../inject#default-implementation)),
  the decorator applies only when config chooses that implementation, so
  every implementation may carry its own `@Middleware` of the same
  interface.

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
  internal/.../*_module.go           # inject modules, one per @Provide/@Invoke/@Config, at the package's path in the project
  client/<client type import path>/<name>_<type>.go # a factory client asked by @Inject (see inject)
  builder/<app>/modules.go           # func Modules() []fx.Option
  builder/<app>/<server>_module.go   # one REST, gRPC or MCP server and its routes/services/tools
  builder/<app>/function_module.go   # the @Function handler and its adapter
  builder/middleware/<package path>/<func>_middleware.go # Decorate<Func>(): the @Middleware decorator
```

`gen/` mirrors the project: a package in `internal/app/api/rest` gets its
module in `gen/internal/app/api/rest`, so apps whose packages live under
`internal/` ([project layout](https://github.com/golang-standards/project-layout))
compile and `gen/` reads like the project. Below `gen/builder/middleware/` and
`gen/implementation/` the leading `internal` is written `internal_`, which Go
needs there (see [inject](../inject#generating)).

`Modules()` returns the doors sorted by name, then the `@Middleware`
decorators the doors reach and, when the app has servers,
`multiserverfx.Module()`, which starts all of them. Each door includes the
inject modules of its handlers, so only what the app serves is built.
Generated files start with `// Code generated by builder; DO NOT EDIT.`; every
run removes the ones no longer produced, including the main of a removed app.

Settings are generated too: a package declares a struct marked `@Config`
whose fields carry `default` and `desc` tags, and inject registers one config
key per field and provides `*Struct` to the constructors that receive it (see
[`inject`](../inject#config)). A door that reaches such a constructor
includes the config module, so the package never imports `wrapper/config`.
Factory clients too: `@Inject(index=0, name=quote, plugins=[otel])` on a
`*resty.Client` (redis, mongo, NATS, Pub/Sub, gocql, AWS) parameter makes
inject build that client with those plugins and its own config under
`boost.factory.resty.clients.quote` (see
[`inject`](../inject#factory-clients)), so the app writes no factory module.

builder owns the main of every app: a `cmd/<app>/main.go` without that header
fails the generation instead of being overwritten. Other files under `cmd/`
(tests next to the main, other commands) are kept.

### Servers

- `rest`, `grpc` and `mcp` use the shipped
  [echo](../fx/modules/factory/contrib/labstack/echo/v4),
  [gRPC](../fx/modules/factory/contrib/google.golang.org/grpc/v1/server) and
  [MCP](../fx/modules/factory/contrib/modelcontextprotocol/go-sdk/v1/server)
  fx modules, configured under `boost.factory.echo`,
  `boost.factory.grpc.server` and `boost.factory.mcp.server`.
- Any other name is a separate server configured under
  `<app>.server.<name>` (`api.server.admin.port`): give it its own port.
- A server name serves one kind of door: REST, gRPC or MCP.
- Plugins are optional and come from `@Provide`: a `[]boostecho.Plugin`
  (REST) or `[]server.Plugin` (gRPC,
  [`grpc/v1/server`](../factory/contrib/google.golang.org/grpc/v1/server);
  MCP, [`go-sdk/v1/server`](../factory/contrib/modelcontextprotocol/go-sdk/v1/server))
  provided with `name=<server>` reaches that server, `rest`, `grpc` and `mcp`
  included; one provided without `name` reaches `rest`, `grpc` or `mcp` too.
  The door includes the provider's inject module, so the generated main needs
  nothing else:

```go
// Plugins tags the admin responses.
// @Provide(name=admin, index=0)
func Plugins() []boostecho.Plugin { ... }
```

## Validation

The command fails, without touching the output directory, when:

- `method` is missing or unknown, or `path` is missing;
- the same `method`+`path` is declared twice in one server, a gRPC
  `register` twice, or an MCP tool name twice;
- an MCP tool name has a character other than letters, digits, `_`, `-` and
  `.`, or more than 128;
- the handler signature does not match the door;
- a method handler has no `@Provide` of its receiver;
- the receiver of a handler method, or a `@GRPC` constructor, is one of the
  implementations of a set (`@Default`): a door is built from one
  constructor, so it takes the chosen value through a type of its own;
- `@GRPC` has no single default `@Provide`, or `register` does not resolve
  to an exported function;
- an app has two `@Function`, or `adapter` is missing or unknown, or its
  list is empty or repeats a name;
- `middlewares` is not written as `[a, b]`, or lists an empty, unknown or
  repeated name;
- a server is used by two kinds of door;
- `rest`, `grpc` or `mcp` gets both a `[]Plugin` with `name=<server>` and
  one without `name`: provide one of them;
- an app or server name is invalid, or an annotation has an unknown key;
- the handler is unexported or in package `main`;
- `@Middleware` lists no plugin, an unknown one or one twice, is repeated,
  is on a struct, is not on a single `@Provide`, decorates a group, a type that is not an
  exported non-generic interface with methods, or a method that breaks the
  rules above, or two constructors decorate the same interface (unless they
  are implementations of the same set);
- `cmd/<app>/main.go` exists and was not generated by builder;
- any inject validation fails (see [`inject`](../inject#validation)).

## Library

```go
// what cmd/builder does: gen/ under the module of src, mains under cmd
err := builder.Generate(ctx, src, "gen", "cmd") // empty cmd: no main

// the files Generate would write, change or remove, relative to src
stale, err := builder.Stale(ctx, src, "gen", "cmd")
```

`Stale` writes nothing, which lets a test of the app check that `gen/` and
the mains are up to date:

```go
func TestGeneratedCodeIsUpToDate(t *testing.T) {
	stale, err := builder.Stale(context.Background(), ".", "gen", "cmd")
	require.NoError(t, err)
	assert.Empty(t, stale, "run go generate")
}
```

The steps underneath, for a tool of its own:

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

## Example

[`examples/simple`](./examples/simple) has two apps:

- `api`: REST `POST /orders` and `GET /orders/:id` on `rest`, `GET`/`HEAD
  /live` on `admin` (built with the `X-Served-By: admin` plugin of
  `probe.Plugins`), the gRPC health service on `grpc` and the `check_stock`
  MCP tool of `stocktool` on `mcp`; `POST /orders` and `check_stock`
  check `stock.Service`, implemented by `stock/memory` (`@Default`, with
  `@Middleware(log, prometheus)`) and `stock/unlimited`, chosen by
  `builder.examples.simple.stock.implementation`
  (`BUILDER_EXAMPLES_SIMPLE_STOCK_IMPLEMENTATION`); `stock/memory` reads its
  sold-out items from its `@Config Options`
  (`builder.examples.simple.stock.memory.soldOut`, env
  `BUILDER_EXAMPLES_SIMPLE_STOCK_MEMORY_SOLD__OUT`, default `sold-out`);
- `worker`: the `ordercreated` function on the cloudevents, the nats or the
  pubsub adapter, chosen by `worker.function.adapter` (`WORKER_FUNCTION_ADAPTER`,
  default `cloudevents`), behind the `recovery` and `logger` middlewares; its
  handler injects the publisher with `drivers=[nats, pubsub]`, so the worker
  links no Kafka client.

`cmd/api/main.go` and `cmd/worker/main.go` are generated. Its tests call the
routes, the gRPC service and the MCP tool of the generated app, check that the handler
receives the stock decorator and that prometheus counted its call, boot the
api with `unlimited` (wired without the decorator), with an unknown name
(the boot fails) and with `SOLD__OUT=book,pen` (the env replaces the
default list), run the function handler, boot the worker with each adapter
(nats on an embedded server) and with an unknown one (the boot fails), run
the generated main with only `WORKER_FUNCTION_ADAPTER=nats` until it records
an order published on NATS and then send it SIGTERM (it exits 0 within a
second), run it with only `WORKER_FUNCTION_ADAPTER=pubsub`, v2 and v1, against
a `pstest` fake until it records a CloudEvent and a raw message, stays up
and exits 0 on SIGTERM, and with no subscription (the boot fails naming it),
build the worker and check with `go version -m` that it has no
`confluent-kafka-go`, boot it with the kafka driver (the boot fails listing
`nats, noop, pubsub`), and fail when `gen/` or a main is stale.

## Non-goals

- A YAML spec: annotations are the only source, and server settings stay in
  boost config.
- Generating OpenAPI or proto from annotations: the spec in `api/` is the
  source, and code may come from it, never the reverse.
- Transport-agnostic handlers: a handler is native to its door.
- Per-route middleware and GraphQL.
- `cache` and `fallback` in `@Middleware` (they need a cache manager or a
  fallback per method), and different middlewares per method.
