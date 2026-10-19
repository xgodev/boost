# Examples

Every example of boost is a builder app: annotated handlers and constructors,
`go generate`, the generated code committed under `gen/`, the main builder
writes under `cmd/<app>/main.go`, and a test that runs the app in process.
Each example is its own Go module next to the package it shows, so the root
`go test ./...` does not reach it; `scripts/examples.sh` vets and tests them
all.

[`boot`](boot) is the smallest one: one handler, one service and the
generated main, whose `boost.Start` loads the config and sets the logger.

[`stock`](stock) is the complete one: a stock with reservation that uses
every resource it can — REST, gRPC and MCP doors, a CloudEvents/NATS
function, service middlewares (log, prometheus), a memory or MongoDB
repository, a resty client, the publisher and `@Config` options, each
switched by config. A resource that lands in boost is added to it in the
same PR. Open `examples/stock` in GoLand and run `cmd/api` and `cmd/worker`
as its `doc.go` says.

## Where they are

| Package | Examples |
|---|---|
| [`bootstrap/function`](../bootstrap/function/examples) | cloudevents http and nats, nats, pubsub, middlewares, recovery |
| [`extra/health`](../extra/health/examples) | simple |
| [`extra/middleware`](../extra/middleware/examples) | simple, service, custom |
| [`extra/middleware/plugins/local/wrapper/cache`](../extra/middleware/plugins/local/wrapper/cache/examples) | service |
| [`factory/local/extra/middleware`](../factory/local/extra/middleware/examples) | simple, error |
| [`factory/contrib/...`](../factory/contrib) | one or more per factory: aws, kafka, elasticsearch, resty, mongo, gocloud pubsub, gocql, grpc, echo, mcp, redis |
| [`wrapper/cache`](../wrapper/cache/examples) | freecache, freecache-binary, multi, plugin |
| [`wrapper/config`](../wrapper/config/examples) | env, env_file |
| [`builder`](../builder/examples) | simple: two apps with rest, grpc, mcp and function doors |
| [`examples/stock`](stock) | the complete app: every resource boost has |
| [`docs/site/example/shop`](../docs/site/example/shop) | the app of the site |

## Shape of an example

The layout is the one of an xgodev service: the handler calls a service, the
service a repository or client, and the domain imports neither transport nor
persistence.

```
doc.go                       # the package doc, how to run it, //go:generate
configs/config.yaml          # when the example needs settings
cmd/<app>/main.go            # generated
cmd/<app>/main_test.go       # runs the app in process
gen/                         # generated, committed
gen_test.go                  # fails when gen/ is stale
internal/app/<app>/<transport>/handler/<entity>/<agent>/handler.go
internal/pkg/domain/model/entity/
internal/pkg/domain/service/<entity>/<agent>/          # interface
internal/pkg/domain/service/<entity>/<agent>/local/    # implementation
internal/pkg/domain/repository/<entity>/               # interface
internal/pkg/datastore/<entity>/<store>/datastore.go   # implementation
```

## Running

```shell
cd wrapper/cache/examples/freecache
go generate ./...
CONF=configs/config.yaml go run ./cmd/api
```

`doc.go` of each example says how to run it and what it needs up (a broker,
a database). The tests run against in-memory fakes or a server the test
starts; the mongo and gocql ones skip unless their env var points at a
database.

```shell
scripts/examples.sh            # every example
scripts/examples.sh wrapper    # the examples under wrapper/
```
