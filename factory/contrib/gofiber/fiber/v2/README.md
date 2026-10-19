# fiber v2

Factory for [Fiber](https://gofiber.io/) **v2** (`github.com/gofiber/fiber/v2`),
the Express-style web framework over [fasthttp](https://github.com/valyala/fasthttp):
handlers are `func(*fiber.Ctx) error` and the error a handler returns goes to
the app error handler. Keys live under `boost.factory.fiber.*`.

## Usage

```go
import (
	"context"

	f "github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2/plugins/extra/error_handler"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2/plugins/local/extra/health"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2/plugins/native/recover"
)

func main() {
	boost.Start()
	ctx := context.Background()

	srv, err := fiber.NewServer(ctx, recover.Register, error_handler.Register, health.Register)
	if err != nil {
		panic(err)
	}

	srv.Get("/hello/:name", func(c *f.Ctx) error {
		return c.SendString("Hello, " + c.Params("name"))
	})

	srv.Serve(ctx) // blocks until srv.Shutdown(ctx)
}
```

Constructors:

| Constructor | Options from |
|---|---|
| `NewServer(ctx, plugins...)` | `boost.factory.fiber.*` |
| `NewServerWithConfigPath(ctx, path, plugins...)` | `boost.factory.fiber.*` overridden by `path.*` |
| `NewServerWithOptions(ctx, opts, plugins...)` | the given `*Options` |

Each returns `(*Server, error)`: a plugin that fails stops the build and its
error is returned. `Serve` blocks and logs any error; `Shutdown` drains the
server within the context deadline, and a `Shutdown` before `Serve` makes
`Serve` return at once. `*Server` is a `multiserver.Server`.

Fiber reads its config only when the app is created, so the app is created
after the plugins: a plugin changes the config through `srv.Config()` (error
handler, JSON encoder) and `srv.App()` is nil while the plugins run. The
middlewares the plugins add with `Use` are added before their routes, so every
plugin middleware wraps every plugin route, in any plugin order. After the
build, `*Server` has the route methods (`Get`, `Post`, `Add`, `All`, `Route`,
`Mount`, `Static`, …) and `App()` returns the `*fiber.App` for the rest.

`ErrorHandlerJSON` and `ErrorHandlerString` answer a `*fiber.Error` with its
code, a `validator.ValidationErrors` with 422 and any other error with the
status its `model/errors` kind maps to; an error ignored with
`errors.IgnoreAsSuccess` answers 200 and a `HEAD` gets no body.

With fx, `fx/modules/factory/contrib/gofiber/fiber/v2.Module()` provides
`*fiber.Server` and `*fiber.App` and joins the multiserver group, so the
multiserver module serves it on the app start and shuts it down on the stop.
Plugins come from the optional `[]fiber.Plugin` in the graph.

## Configuration

| Key | Default | Meaning |
|---|---|---|
| `boost.factory.fiber.port` | `8082` | listen port |
| `boost.factory.fiber.type` | `REST` | `REST` makes the error handler plugin answer JSON |
| `boost.factory.fiber.config.network` | `tcp4` | `tcp`, `tcp4` or `tcp6` |
| `boost.factory.fiber.config.prefork` | `false` | one process per CPU on the same port |
| `boost.factory.fiber.config.bodyLimit` | `4194304` | max request body, in bytes |
| `boost.factory.fiber.config.concurrency` | `262144` | max concurrent connections |
| `boost.factory.fiber.config.readTimeout` / `writeTimeout` / `idleTimeout` | `0s` | `0s` is unlimited |
| `boost.factory.fiber.config.readBufferSize` / `writeBufferSize` | `4096` | per-connection buffers |
| `boost.factory.fiber.config.disableStartupMessage` | `true` | hide the fiber banner; boost logs the start |
| `boost.factory.fiber.config.proxyHeader` | `""` | header the client IP is read from |
| `boost.factory.fiber.config.enableTrustedProxyCheck` / `trustedProxies` | `false` / `[]` | trust the proxy headers only from these |
| `boost.factory.fiber.config.serverHeader`, `strictRouting`, `caseSensitive`, `immutable`, `unescapePath`, `etag`, `compressedFileSuffix`, `getOnly`, `disableKeepalive`, `disableDefaultDate`, `disableDefaultContentType`, `disableHeaderNormalizing`, `reduceMemoryUsage` | fiber's | the `fiber.Config` field of the same name |

## Plugins

Every plugin has a package `Register(ctx, srv) error` that reads its keys
under `boost.factory.fiber.plugins.<name>`, and the constructor trio
`NewX()`, `NewXWithConfigPath(path)`, `NewXWithOptions(opts)`.

| Plugin | Package | Keys (`…plugins.`) |
|---|---|---|
| CORS | `plugins/native/cors` | `cors.enabled`, `cors.allowed.{origins,headers,methods,credentials}`, `cors.exposed.headers`, `cors.maxAge` |
| Compress | `plugins/native/compress` | `compress.enabled`, `compress.level` (`-1`…`2`) |
| ETag | `plugins/native/etag` | `etag.enabled`, `etag.weak` |
| Monitor | `plugins/native/monitor` | `monitor.enabled` (`false`), `monitor.route` (`/monitor`), `monitor.title` |
| pprof | `plugins/native/pprof` | `pprof.enabled` (`false`), `pprof.prefix` |
| Recover | `plugins/native/recover` | `recover.enabled`, `recover.stackTrace` |
| Request ID | `plugins/native/requestid` | `requestId.enabled`, `requestId.header` (`X-Request-ID`) |
| Error handler | `plugins/extra/error_handler` | `errorHandler.enabled` |
| Health | `plugins/local/extra/health` | `health.enabled`, `health.route` (`/health`) |
| Multiserver check | `plugins/local/extra/multi-server` | `multiServer.enabled`, `multiServer.route` (`/check`) |
| Resource status | `plugins/local/model/restresponse` | `status.enabled`, `status.route` (`/resource-status`) |
| Request log | `plugins/local/wrapper/log` | `logger.enabled`, `logger.level` (`DEBUG`) |
| Prometheus | `plugins/contrib/prometheus/client_golang/v1` | `prometheus.enabled`, `prometheus.route.{path,enabled}`, `prometheus.collector.enabled` |
| Datadog | `plugins/contrib/datadog/dd-trace-go/v2` | `datadog.enabled` (no-op while the tracer is off) |
| goccy/go-json | `plugins/contrib/goccy/go-json/v0` | `goccyGoJson.enabled` |
| bytedance/sonic | `plugins/contrib/bytedance/sonic/v1` | `bytedanceSonic.enabled` |

The request log and Prometheus middlewares pass a handler error through the
app error handler, so they see the status the client gets. The Prometheus
collectors (`fiber_requests_total`, `fiber_request_duration_seconds`, labels
`code`, `method`, `url`) register once per process in the default registry,
so every fiber server of the app reports into the same metrics. `url` is the
route path (`/orders/:id`), one series per route. pprof and monitor expose
internals and are off by default.
