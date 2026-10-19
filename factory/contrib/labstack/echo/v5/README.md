# echo v5

Factory for [Echo](https://echo.labstack.com/) **v5** (`github.com/labstack/echo/v5`).
It lives beside [`echo/v4`](../v4/README.md): same options, same config keys
(`boost.factory.echo.*`) and the same plugin set, written against the v5 API.
An app picks one of the two; a service moves from v4 to v5 by changing the
import paths and the handler signature.

## Usage

```go
import (
	"context"
	"net/http"

	e "github.com/labstack/echo/v5"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5/plugins/local/extra/health"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5/plugins/native/recover"
)

func main() {
	boost.Start()
	ctx := context.Background()

	srv, err := echo.NewServer(ctx, recover.Register, health.Register)
	if err != nil {
		panic(err)
	}

	srv.GET("/hello", func(c *e.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	srv.Serve(ctx) // blocks until srv.Shutdown(ctx)
}
```

Constructors:

| Constructor | Options from |
|---|---|
| `NewServer(ctx, plugins...)` | `boost.factory.echo.*` |
| `NewServerWithConfigPath(ctx, path, plugins...)` | `boost.factory.echo.*` overridden by `path.*` |
| `NewServerWithOptions(ctx, opts, plugins...)` | the given `*Options` |

Each returns `(*Server, error)`: a plugin that fails stops the build and its
error is returned. `Serve` blocks and logs any error but the
`http.ErrServerClosed` that `Shutdown` causes; `Shutdown` drains the server
within the context deadline. `*Server` is a `multiserver.Server`.

With fx, `fx/modules/factory/contrib/labstack/echo/v5.Module()` provides
`*echo.Server` and `*e.Echo` and joins the multiserver group, so the
multiserver module serves it on the app start and shuts it down on the stop.
Plugins come from the optional `[]echo.Plugin` in the graph.

## Configuration

| Key | Default | Meaning |
|---|---|---|
| `boost.factory.echo.port` | `8080` | listen port |
| `boost.factory.echo.type` | `REST` | `REST` answers errors as JSON, anything else as text |
| `boost.factory.echo.protocol` | `HTTP` | `HTTP` or `H2C` |
| `boost.factory.echo.hideBanner` | `true` | hide the `Echo v5.x` line |
| `boost.factory.echo.disableHTTP2` | `false` | TLS offers only HTTP/1.1 |
| `boost.factory.echo.tls.enabled` | `false` | serve TLS |
| `boost.factory.echo.tls.type` | `AUTO` | `AUTO` (self-signed) or `FILE` |
| `boost.factory.echo.tls.auto.host` | `localhost` | host of the self-signed certificate |
| `boost.factory.echo.tls.file.cert` / `.key` | `""` | PEM files for `FILE` |
| `boost.factory.echo.hc2.*` | | HTTP/2 server options for `H2C` |

## Plugins

Every plugin has a package `Register(ctx, srv) error` that reads its keys
under `boost.factory.echo.plugins.<name>`, and the constructor trio
`NewX()`, `NewXWithConfigPath(path)`, `NewXWithOptions(opts)`.

| Plugin | Package | Keys (`…plugins.`) |
|---|---|---|
| CORS | `plugins/native/cors` | `cors.enabled`, `cors.allowed.{origins,headers,methods,credentials}`, `cors.exposed.headers`, `cors.maxAge` |
| Gzip | `plugins/native/gzip` | `gzip.enabled`, `gzip.level` |
| Recover | `plugins/native/recover` | `recover.enabled` |
| Request ID | `plugins/native/requestid` | `requestId.enabled` |
| Body limit | `plugins/native/bodylimit` | `bodyLimit.enabled`, `bodyLimit.size` (`8M`) |
| Body dump | `plugins/native/bodydump` | `bodydump.enabled` |
| Error handler | `plugins/extra/error_handler` | `errorHandler.enabled` |
| Semaphore | `plugins/extra/semaphore` | `semaphore.enabled`, `semaphore.limit` (`10000`, then 503) |
| Health | `plugins/local/extra/health` | `health.enabled`, `health.route` (`/health`) |
| Multiserver check | `plugins/local/extra/multi-server` | `multiServer.enabled`, `multiServer.route` (`/check`) |
| Resource status | `plugins/local/model/restresponse` | `status.enabled`, `status.route` (`/resource-status`) |
| Request log | `plugins/local/wrapper/log` | `logger.enabled`, `logger.level` (`DEBUG`) |
| Prometheus | `plugins/contrib/labstack/echo-contrib/v5/echoprometheus` | `prometheus.enabled`, `prometheus.route.{path,enabled}`, `prometheus.collector.enabled` |
| pprof | `plugins/contrib/labstack/echo-contrib/v5/pprof` | `pprof.enabled`, `pprof.prefix` (`/debug/pprof`) |
| OpenTelemetry | `plugins/contrib/labstack/echo-otel/v5` | `otel.enabled`, `otel.serverName` |
| Datadog | `plugins/contrib/datadog/dd-trace-go/v2` | `datadog.enabled` (no-op while the tracer is off) |
| goccy/go-json | `plugins/contrib/goccy/go-json/v0` | `goccyGoJson.enabled` |
| bytedance/sonic | `plugins/contrib/bytedance/sonic/v1` | `bytedanceSonic.enabled` |
| Swagger UI | [`plugins/contrib/swaggo/echo-swagger/v2`](plugins/contrib/swaggo/echo-swagger/v2/README.md) | `swagger.enabled`, `swagger.route` (`/swagger/*`) |

The Prometheus collectors register once per process in the default registry,
so every echo server of the app reports into the same `echo_*` metrics.
The OpenTelemetry plugin uses the providers of
`factory/contrib/go.opentelemetry.io/otel/v1`.

## Differences from v4

- Handlers take `*echo.Context` (v5) instead of `echo.Context`; route methods
  return `echo.RouteInfo`; `ErrorHandlerJSON`/`ErrorHandlerString` take
  `(c, err)`.
- Echo logs through `log/slog`; `WrapLogger` bridges it to the boost logger.
- Constructors and plugins return their errors instead of panicking or
  calling `log.Fatal`.
- `cors.allowed.credentials` defaults to `false`: echo v5 refuses credentials
  with the `*` origin. Set explicit origins to turn it on.
- `bodyLimit.*` and `bodydump.*` keys are registered under their plugin path
  (v4 missed them) and H2C reads `hc2.*` under the server path.
- Prometheus, pprof and OpenTelemetry use the v5 `echo-contrib`/`echo-otel`
  packages; Swagger UI uses `swaggo/echo-swagger/v2` (v4 has its v1 plugin).
- The [builder](../../../../../builder/README.md) REST door still generates
  echo v4 code.
