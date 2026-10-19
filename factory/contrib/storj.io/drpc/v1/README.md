# drpc

Factories for [storj drpc](https://github.com/storj/drpc) v1, a lightweight
drop-in replacement for gRPC: a server that joins `extra/multiserver` and a
client connection.

## Server

```go
srv, err := server.NewServer(ctx, func(_ context.Context, s *server.Server) error {
	return pb.DRPCRegisterOrders(s.Mux(), ordersImpl)
})
go srv.Serve(ctx)
defer srv.Shutdown(ctx)
```

`Register(impl, desc)` registers a service on the mux directly. `Serve` blocks
until `Shutdown` (or `ctx`) stops it; `Shutdown` closes the listener and the
open connections and waits for `Serve`, bounded by its `ctx`. A `Shutdown`
before `Serve` makes `Serve` return at once.

| Key | Default | Description |
|---|---|---|
| `boost.factory.drpc.server.port` | `9091` | server drpc port |
| `boost.factory.drpc.server.collectStats` | `false` | collect per rpc stats (`Server().Stats()`) |

## Client

```go
conn, err := client.NewClientConn(ctx)
defer conn.Close()
orders := pb.NewDRPCOrdersClient(conn)
```

The connection is dialed when it is built and the caller closes it. A plugin
that fails closes it and returns the error.

| Key | Default | Description |
|---|---|---|
| `boost.factory.drpc.client.host` | `localhost` | defines host |
| `boost.factory.drpc.client.port` | `9091` | defines port |
| `boost.factory.drpc.client.timeout` | `10s` | defines dial timeout |

`NewServerWithConfigPath` / `NewClientConnWithConfigPath` read the same keys
under another path.

## Plugins

`server.Plugin = factory.Plugin[*server.Server]` and
`client.Plugin = factory.Plugin[*drpcconn.Conn]`.

## fx

`fx/modules/factory/contrib/storj.io/drpc/v1/server` provides `*server.Server`
and `*drpcmux.Mux` and joins the multiserver group.
`fx/modules/factory/contrib/storj.io/drpc/v1/client` provides `*drpcconn.Conn`
and closes it on the app stop.
