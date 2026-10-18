# MCP client for Boost

Factory for an [MCP](https://modelcontextprotocol.io) client built on the
official SDK [`github.com/modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk).
The constructor connects: it returns a client with an open session over
streamable HTTP or a subprocess speaking stdio.

```go
import "github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
```

## Usage

```go
boost.Start()

cli, err := client.NewClient(ctx, mcplog.Register, otel.Register)
if err != nil {
    log.Fatal(err)
}
defer cli.Close()

res, err := cli.Session().CallTool(ctx, &mcp.CallToolParams{
    Name:      "greet",
    Arguments: map[string]any{"name": "boost"},
})
```

- `Session()` is the SDK `*mcp.ClientSession`: tools, prompts, resources.
- `Instance()` is the SDK `*mcp.Client`.
- Plugins run before the connection, so they see `initialize` too.
- `NewClientWithConfigPath(ctx, path)` reads another upstream from its own root
  (e.g. `boost.factory.mcp.client.github`), falling back to the defaults below.

Full example: [`examples/greet`](examples/greet/main.go), against
[`server/examples/greet`](../server/examples/greet/main.go).

## Configuration

Root `boost.factory.mcp.client` (env `BOOST_FACTORY_MCP_CLIENT_*`).

| Key | Default | Description |
|---|---|---|
| `.name` | `""` | implementation name sent on initialize; empty uses `boost.applicationName` |
| `.version` | `""` | implementation version sent on initialize; empty uses the binary's main module version (`go version -m`), `0.0.0` when unstamped |
| `.transport` | `HTTP` | `HTTP` (streamable HTTP endpoint) or `COMMAND` (subprocess over stdio) |
| `.endpoint` | `http://localhost:8090/mcp` | HTTP transport server URL |
| `.maxRetries` | `5` | HTTP reconnect attempts; negative disables (env `..._MAX__RETRIES`) |
| `.command` | `""` | COMMAND transport executable |
| `.args` | `[]` | COMMAND transport arguments |
| `.keepAlive` | `0` | ping interval that closes a dead session; 0 disables (env `..._KEEP__ALIVE`) |

## Plugins

Each plugin is a `client.Plugin` with its own `.enabled` key under
`boost.factory.mcp.client.plugins`.

| Plugin | Keys | What it does |
|---|---|---|
| [`plugins/local/wrapper/log`](plugins/local/wrapper/log) | `.log.enabled` (`true`), `.log.level` (`DEBUG`) | logs method, tool, duration and error of every call through `wrapper/log` |
| [`plugins/contrib/go.opentelemetry.io/otel/v1`](plugins/contrib/go.opentelemetry.io/otel/v1) | `.otel.enabled` (`true`) | one client span per call from the boost otel factory, injecting the trace context into `_meta` |

## fx

`fx/modules/factory/contrib/modelcontextprotocol/go-sdk/v1/client` provides
`*client.Client` and `*mcp.ClientSession`, takes `[]client.Plugin` from the
graph and closes the session on app stop.
