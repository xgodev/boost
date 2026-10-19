# fasthttp client v1

Factory for a [fasthttp](https://github.com/valyala/fasthttp) `*fasthttp.Client`
configured from `boost.factory.fasthttp.client.*`.

## Usage

```go
import (
	"context"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/valyala/fasthttp/v1/client"
)

func main() {
	boost.Start()
	ctx := context.Background()

	c, err := client.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	status, body, err := c.Get(nil, "http://localhost:8082/health")
	// ...
}
```

| Constructor | Options from |
|---|---|
| `NewClient(ctx, plugins...)` | `boost.factory.fasthttp.client.*` |
| `NewClientWithConfigPath(ctx, path, plugins...)` | `boost.factory.fasthttp.client.*` overridden by `path.*` |
| `NewClientWithOptions(ctx, opts, plugins...)` | the given `*Options` |

Each returns `(*fasthttp.Client, error)`. A `Plugin` is
`factory.Plugin[*fasthttp.Client]`: it changes the client before it is
returned, and one that fails stops the build.

With fx, `fx/modules/factory/contrib/valyala/fasthttp/v1/client.Module()`
provides `*fasthttp.Client`, with the optional `[]client.Plugin` of the graph.

## Configuration

| Key (`boost.factory.fasthttp.client.`) | Default | Meaning |
|---|---|---|
| `name` | `fasthttp-client` | User-Agent of the requests |
| `noDefaultUserAgentHeader` | `false` | send no User-Agent |
| `maxConnsPerHost` | `512` | max concurrent connections per host |
| `maxConnWaitTimeout` | `0s` | wait for a free connection; `0s` fails at once |
| `readTimeout` / `writeTimeout` | `0s` | full response read / request write; `0s` is unlimited |
| `readBufferSize` / `writeBufferSize` | `0` | per-connection buffers; `0` is the fasthttp default |
| `maxIdleConnDuration` | `10s` | idle keep-alive connections close after it |
| `maxConnDuration` | `0s` | keep-alive connections close after it; `0s` is unlimited |
| `maxResponseBodySize` | `52428800` | max response body, in bytes |
| `maxIdemponentCallAttempts` | `5` | attempts for idempotent calls |
| `disableHeaderNamesNormalizing` | `false` | send header names as they are |
| `dialDualStack` | `false` | try both IPv4 and IPv6 |
