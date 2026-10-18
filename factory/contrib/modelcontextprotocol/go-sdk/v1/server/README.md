# MCP server for Boost

Factory for an [MCP](https://modelcontextprotocol.io) server built on the
official SDK [`github.com/modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk).
It serves over streamable HTTP or stdio and satisfies `extra/multiserver.Server`
(`Serve(ctx)` / `Shutdown(ctx)`).

```go
import "github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
```

## Usage

```go
boost.Start()

srv, err := server.NewServer(ctx, mcplog.Register, otel.Register)
if err != nil {
    log.Fatal(err)
}
mcp.AddTool(srv.Instance(), &mcp.Tool{Name: "greet", Description: "says hello"}, greet)

srv.Serve(ctx) // HTTP: blocks until Shutdown
```

- `Instance()` is the SDK `*mcp.Server`: register tools, prompts and resources on it.
- `Handler()` is the streamable HTTP handler, to mount on your own mux (echo, net/http).
- Plugins run once, inside the constructor; a plugin error fails the constructor.
- `NewServerWithConfigPath(ctx, path)` reads a second server from another root
  (e.g. `boost.factory.mcp.admin`), falling back to the defaults below.

Full example: [`examples/greet`](examples/greet/main.go).

## Configuration

Root `boost.factory.mcp.server` (env `BOOST_FACTORY_MCP_SERVER_*`).

| Key | Default | Description |
|---|---|---|
| `.name` | `""` | implementation name sent on initialize; empty uses `boost.applicationName` |
| `.version` | `1.0.0` | implementation version sent on initialize |
| `.instructions` | `""` | instructions the client may add to its system prompt |
| `.transport` | `HTTP` | `HTTP` (streamable HTTP on port/path) or `STDIO` |
| `.port` | `8090` | HTTP transport port |
| `.path` | `/mcp` | HTTP transport path |
| `.stateless` | `false` | HTTP without session state (no server-to-client requests) |
| `.jsonResponse` | `false` | HTTP answers `application/json` instead of `text/event-stream` (env `..._JSON__RESPONSE`) |
| `.sessionTimeout` | `0` | HTTP idle session timeout; 0 never closes (env `..._SESSION__TIMEOUT`) |
| `.keepAlive` | `0` | ping interval that closes dead sessions; 0 disables (env `..._KEEP__ALIVE`) |

## Plugins

Each plugin is a `server.Plugin` with its own `.enabled` key under
`boost.factory.mcp.server.plugins`.

| Plugin | Keys | What it does |
|---|---|---|
| [`plugins/local/wrapper/log`](plugins/local/wrapper/log) | `.log.enabled` (`true`), `.log.level` (`DEBUG`) | logs method, tool, duration and error of every call through `wrapper/log`; errors log at ERROR, tool results with `isError` are flagged |
| [`plugins/contrib/go.opentelemetry.io/otel/v1`](plugins/contrib/go.opentelemetry.io/otel/v1) | `.otel.enabled` (`true`) | one server span per call from the boost otel factory, continuing the trace the client sent in `_meta` |

## fx

`fx/modules/factory/contrib/modelcontextprotocol/go-sdk/v1/server` provides
`*server.Server` and `*mcp.Server`, takes `[]server.Plugin` from the graph and
joins the `multiserver` group.
