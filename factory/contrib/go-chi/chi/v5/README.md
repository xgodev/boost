# chi v5

Factory for [chi](https://go-chi.io/) **v5** (`github.com/go-chi/chi/v5`), the
`net/http` router: handlers are plain `http.HandlerFunc` and middlewares are
`func(http.Handler) http.Handler`, so any `net/http` middleware plugs in.
Same server options as the [echo](../../../labstack/echo/v5/README.md) factory
(port, `HTTP`/`H2C`, TLS), under `boost.factory.chi.*`.

## Usage

```go
import (
	"context"
	"net/http"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/local/extra/health"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/native/recoverer"
)

func main() {
	boost.Start()
	ctx := context.Background()

	srv, err := chi.NewServer(ctx, recoverer.Register, health.Register)
	if err != nil {
		panic(err)
	}

	srv.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, " + r.PathValue("name")))
	})

	srv.Serve(ctx) // blocks until srv.Shutdown(ctx)
}
```

Constructors:

| Constructor | Options from |
|---|---|
| `NewServer(ctx, plugins...)` | `boost.factory.chi.*` |
| `NewServerWithConfigPath(ctx, path, plugins...)` | `boost.factory.chi.*` overridden by `path.*` |
| `NewServerWithOptions(ctx, opts, plugins...)` | the given `*Options` |

Each returns `(*Server, error)`: a plugin that fails stops the build and its
error is returned. `Serve` blocks and logs any error but the
`http.ErrServerClosed` that `Shutdown` causes; `Shutdown` drains the server
within the context deadline. `*Server` is a `multiserver.Server`.

`*Server` has the `chi.Router` route methods (`Get`, `Post`, `Route`, `Group`,
`Mount`, `NotFound`, …) and `Mux()` returns the `*chi.Mux`. chi panics on a
middleware added after a route, so while the plugins run their routes are held
and added after the last plugin: every plugin middleware wraps every plugin
route, in any plugin order. Plugins call `Use`, never `With`.

With fx, `fx/modules/factory/contrib/go-chi/chi/v5.Module()` provides
`*chi.Server` and `*chi.Mux` and joins the multiserver group, so the
multiserver module serves it on the app start and shuts it down on the stop.
Plugins come from the optional `[]chi.Plugin` in the graph.

## Configuration

| Key | Default | Meaning |
|---|---|---|
| `boost.factory.chi.port` | `8080` | listen port |
| `boost.factory.chi.protocol` | `HTTP` | `HTTP` or `H2C` |
| `boost.factory.chi.disableHTTP2` | `false` | TLS offers only HTTP/1.1 |
| `boost.factory.chi.tls.enabled` | `false` | serve TLS |
| `boost.factory.chi.tls.type` | `AUTO` | `AUTO` (self-signed) or `FILE` |
| `boost.factory.chi.tls.auto.host` | `localhost` | host of the self-signed certificate |
| `boost.factory.chi.tls.file.cert` / `.key` | `""` | PEM files for `FILE` |
| `boost.factory.chi.hc2.*` | | HTTP/2 server options for `H2C` |

## Plugins

Every plugin has a package `Register(ctx, srv) error` that reads its keys
under `boost.factory.chi.plugins.<name>`, and the constructor trio
`NewX()`, `NewXWithConfigPath(path)`, `NewXWithOptions(opts)`.

| Plugin | Package | Keys (`…plugins.`) |
|---|---|---|
| CORS | `plugins/native/cors` | `cors.enabled`, `cors.allowed.{origins,headers,methods,credentials}`, `cors.exposed.headers`, `cors.maxAge` |
| Real IP | `plugins/native/realip` | `realip.enabled` (trusts `X-Real-IP`/`X-Forwarded-For`: only behind a proxy) |
| Recoverer | `plugins/native/recoverer` | `recoverer.enabled` |
| Request ID | `plugins/native/requestid` | `requestid.enabled` (also answers `X-Request-Id`) |
| Strip slashes | `plugins/native/stripslashes` | `stripslashes.enabled` |
| Health | `plugins/local/extra/health` | `health.enabled`, `health.route` (`/health`) |
| Multiserver check | `plugins/local/extra/multi-server` | `multiserver.enabled`, `multiserver.route` (`/check`) |
| Resource status | `plugins/local/model/restresponse` | `status.enabled`, `status.route` (`/resource-status`) |
| Request log | `plugins/local/wrapper/log` | `logger.enabled`, `logger.level` (`DEBUG`) |
| Prometheus | `plugins/contrib/prometheus/client_golang/v1` | `prometheus.enabled`, `prometheus.route.{path,enabled}`, `prometheus.collector.enabled` |
| OpenTelemetry | `plugins/contrib/go.opentelemetry.io/otelhttp/v0` | `otel.enabled`, `otel.serverName` |
| Datadog | `plugins/contrib/datadog/dd-trace-go/v2` | `datadog.enabled` (no-op while the tracer is off) |

The Prometheus collectors (`chi_requests_total`, `chi_request_duration_seconds`,
labels `code`, `method`, `url`) register once per process in the default
registry, so every chi server of the app reports into the same metrics. `url`
is the route pattern (`/orders/{id}`), one series per route. OpenTelemetry
names spans `METHOD /route/{param}` and uses the providers of
`factory/contrib/go.opentelemetry.io/otel/v1`.
