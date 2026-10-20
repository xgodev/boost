# stock

The complete example app of boost: a small stock with reservation, built from
every boost resource it can use. It is a living example: a resource that lands
in boost is added here in the same PR.

| App | What it does |
|---|---|
| `api` | create, find and reserve stock on REST (`:8080`), gRPC (`:9090`) and MCP (`:8090/mcp`) |
| `worker` | reserves the item of each `order.created` CloudEvent it receives, over HTTP (`:8081`) or NATS |
| `stock-check` | a job: logs the units of the skus it is given and exits, 1 when a sku is not in stock |

A reservation is announced as a `stock.reserved` event on boost's publisher.
The defaults need nothing running: the stock is in memory, seeded with
`tshirt=3` and `mug=1`, and the publisher is `noop`.

## Run

From this directory, each app in its own terminal:

```shell
CONF=configs/api.yaml go run ./cmd/api
CONF=configs/worker.yaml go run ./cmd/worker
INTERNAL_APP_STOCKCHECK_JOB_STOCK_CHECKER_SKUS=tshirt,mug go run ./cmd/stock-check
```

`configs/worker.yaml` puts the worker's HTTP door on `8081` and makes it
reserve by calling app `api`, so run `api` first.

In GoLand: open `examples/stock`, add a Go Build configuration per app (kind
Package, `github.com/xgodev/boost/examples/stock/cmd/api`), working directory
`examples/stock`, and the `CONF` variable above.

## Call it

In this order, on a fresh `api` and `worker`.

REST, on app `api`:

```console
$ curl localhost:8080/stocks/tshirt
{"sku":"tshirt","units":3}

$ curl -XPOST localhost:8080/stocks/tshirt/reserve
{"sku":"tshirt","units":2}

$ curl localhost:8080/stocks/nope
{"httpStatusCode":404,"errorCode":"404","message":"nope is not in stock"}

$ curl -XPOST localhost:8080/stocks -H 'Content-Type: application/json' -d '{"sku":"cap","units":1}'
{"sku":"cap","units":1}

$ curl -XPOST localhost:8080/stocks/cap/reserve
{"sku":"cap","units":0}

$ curl -XPOST localhost:8080/stocks/cap/reserve
{"httpStatusCode":409,"errorCode":"409","message":"cap is sold out"}

$ curl localhost:8080/health
{"status":"OK"}
```

The create answers 201. `/metrics` serves Prometheus, and every response
carries an `X-Request-Id`.

MCP: point any MCP client at `http://localhost:8090/mcp` (Claude Code:
`claude mcp add --transport http stock http://localhost:8090/mcp`). The tool
`find_stock` with `{"sku":"mug"}` answers `{"sku":"mug","units":1}`.

gRPC (`api/proto/stock/v1/stock.proto`; the server has reflection on):

```console
$ grpcurl -plaintext -d '{"sku":"mug"}' localhost:9090 stock.v1.StockReserver/Reserve
{
  "sku": "mug"
}
$ grpcurl -plaintext -d '{"sku":"mug"}' localhost:9090 stock.v1.StockReserver/Reserve
ERROR:
  Code: Aborted
  Message: mug is sold out
```

`units` is 0, so the JSON leaves it out.

The worker, with one `order.created` CloudEvent:

```console
$ curl localhost:8081 -H 'Ce-Id: 1' -H 'Ce-Source: shop' -H 'Ce-Type: order.created' \
    -H 'Ce-Specversion: 1.0' -H 'Content-Type: application/json' -d '{"sku":"tshirt"}'
```

Its log says `reserved one tshirt for event 1, 1 left`, and app `api` now
answers `{"sku":"tshirt","units":1}`.

The job:

```console
$ INTERNAL_APP_STOCKCHECK_JOB_STOCK_CHECKER_SKUS=tshirt,nope go run ./cmd/stock-check
INF tshirt has 3 units
ERR job failed: checking nope: nope is not in stock          # exit 1
```

## Config

Every switch is a config key; the env var is the key in upper case with `_`.

| Env var | Effect |
|---|---|
| `INTERNAL_PKG_DOMAIN_REPOSITORY_STOCK_IMPLEMENTATION=mongodb` | stock in MongoDB (`BOOST_FACTORY_MONGO_URI`) |
| `INTERNAL_PKG_DOMAIN_REPOSITORY_STOCK_IMPLEMENTATION=oracle` | stock in Oracle (`BOOST_FACTORY_GODROR_CONNECTSTRING`, `_USERNAME`, `_PASSWORD`); builds with cgo |
| `INTERNAL_PKG_DOMAIN_SERVICE_STOCK_RESERVER_IMPLEMENTATION=resty` | reserve by calling app `api` (`BOOST_FACTORY_RESTY_CLIENTS_STOCK_HOST`) |
| `BOOST_WRAPPER_PUBLISHER_DRIVER_NAME=nats` | announce `stock.reserved` on NATS (`BOOST_FACTORY_NATS_URL`) |
| `WORKER_FUNCTION_ADAPTER=nats` | the worker reads `order.created` from NATS (`BOOST_BOOTSTRAP_FUNCTION_ADAPTER_NATS_SUBJECTS`) |
| `BOOST_FACTORY_ECHO_PLUGINS_PROMETHEUS_ENABLED=false` | turns one plugin of a server off; the others stay |
| `INTERNAL_APP_STOCKCHECK_JOB_STOCK_CHECKER_SKUS=tshirt,mug` | the skus the job checks |

## Where each resource is

| Resource | File | Annotation |
|---|---|---|
| REST door | `internal/app/api/rest/handler/stock/{creator,finder,reserver}/handler.go` | `@REST` |
| gRPC door | `internal/app/api/grpc/handler/stock/{finder,reserver}/handler.go` | `@GRPC` |
| MCP tool | `internal/app/api/mcp/handler/stock/finder/handler.go` | `@MCP` |
| CloudEvents / NATS function | `internal/app/worker/cloudevents/handler/stock/reserver/handler.go` | `@Function(adapter=[cloudevents, nats])` |
| Job | `internal/app/stockcheck/job/stock/checker/handler.go` | `@Job` |
| Server plugins | `internal/pkg/server/{rest,grpc,mcp,cloudevents}/doc.go` | `@Server` |
| REST model | `internal/pkg/server/rest/model/stock.go` | -- |
| Service middlewares (log, prometheus) | `internal/pkg/domain/service/stock/*/local/service.go` | `@Middleware` |
| Implementation by config | `internal/pkg/datastore/stock/memory/datastore.go`, `.../reserver/local/service.go` | `@Default` |
| MongoDB datastore | `internal/pkg/datastore/stock/mongodb/datastore.go` | `@Inject(plugins=[...])` |
| Oracle datastore (`*sql.DB`) | `internal/pkg/datastore/stock/oracle/datastore.go` | `@Inject(drivers=[godror])` |
| resty client | `internal/pkg/client/stock/reserver/resty/client.go` | `@Inject(name=stock)` |
| Publisher | `internal/pkg/publisher/stock/reserved/boost/publisher.go` | `@Inject(drivers=[nats])` |
| Options | the `options.go` next to each of them | `@Config` |

The domain (`internal/pkg/domain`) imports neither transport nor persistence:
the handler calls a service, the service a repository or client.

## Generate and test

```shell
go generate ./...
go test ./...
```

`gen/` and `cmd/*/main.go` are generated and committed; `gen_test.go` fails
when they are stale. `api/proto/stock/v1` is regenerated with the `protoc`
line of its `doc.go`. The tests run each app in process on the defaults and
on each switch above; the MongoDB and Oracle ones skip unless
`BOOST_FACTORY_MONGO_URI` or `BOOST_FACTORY_GODROR_CONNECTSTRING` points at a
database.
