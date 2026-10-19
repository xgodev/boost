# Google Cloud Bigtable Factory (v1)

## Overview

Builds a [`cloud.google.com/go/bigtable`](https://pkg.go.dev/cloud.google.com/go/bigtable)
data client from Boost config: project, instance and credentials.

## Features

- **Client from config**: `NewClient` reads `boost.factory.bigtable.*`
- **Custom config path or options**: `NewClientWithConfigPath`, `NewClientWithOptions`
- **Named clients**: `ConfigAdd(path)` registers the keys under any path, so one app can hold more than one client
- **gRPC plugins**: any `clientgrpc.Plugin` adds dial options
- **Emulator**: `BIGTABLE_EMULATOR_HOST` points the client at a local emulator
- **fx module**: `fx/modules/factory/contrib/cloud.google.com/bigtable/v1` provides `*bigtable.Client` and closes it on app stop
- **Logging**: the client creation is logged through `wrapper/log`

## Usage

### Basic Usage

```go
package main

import (
    "context"

    gbigtable "cloud.google.com/go/bigtable"
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/cloud.google.com/bigtable/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    boost.Start()

    ctx := context.Background()

    client, err := bigtable.NewClient(ctx)
    if err != nil {
        log.Fatalf("failed to create Bigtable client: %v", err)
    }
    defer client.Close()

    orders := client.Open("orders")

    mut := gbigtable.NewMutation()
    mut.Set("cf", "status", gbigtable.Now(), []byte("paid"))
    if err := orders.Apply(ctx, "order#1", mut); err != nil {
        log.Fatalf("failed to write row: %v", err)
    }

    row, err := orders.ReadRow(ctx, "order#1")
    if err != nil {
        log.Fatalf("failed to read row: %v", err)
    }
    log.Infof("status: %s", row["cf"][0].Value)
}
```

`project` and `instance` are required: without them the constructors return a
`NotValid` error before dialing, instead of the SDK failing on the first call.

### Named Clients

```go
bigtable.ConfigAdd("boost.factory.bigtable.clients.orders")

client, err := bigtable.NewClientWithConfigPath(ctx, "boost.factory.bigtable.clients.orders")
```

```yaml
boost:
  factory:
    bigtable:
      clients:
        orders:
          project: my-gcp-project
          instance: orders-instance
```

`ConfigAdd` registers the keys under the path with empty defaults, so the
named client sets its own `project` and `instance`.

### Custom Options

```go
options, err := bigtable.NewOptions()
if err != nil {
    log.Fatalf("failed to load options: %v", err)
}

options.Instance = "orders-instance"

client, err := bigtable.NewClientWithOptions(ctx, options)
```

### Using Plugins

```go
import (
    otelgrpc "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client/plugins/contrib/go.opentelemetry.io/contrib/v0"
)

client, err := bigtable.NewClient(ctx, otelgrpc.Register)
```

A plugin is a `clientgrpc.Plugin`:
`func(ctx context.Context) ([]grpc.DialOption, []grpc.CallOption)`. The
client uses the dial options; the Bigtable SDK has no hook for call options,
so a plugin that only returns call options (like `native/compressor`) has no
effect.

### fx Module

```go
import (
    btfx "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/bigtable/v1"
    "go.uber.org/fx"
)

fx.New(
    btfx.Module(),
    fx.Invoke(func(c *gbigtable.Client) { /* ... */ }),
)
```

Plugins in the graph as `[]clientgrpc.Plugin` are passed to the client. The
client is closed on app stop.

## Emulator

Run the emulator from the gcloud CLI image and point the app at it:

```bash
docker run --rm -p 8086:8086 gcr.io/google.com/cloudsdktool/google-cloud-cli:emulators \
  gcloud beta emulators bigtable start --host-port=0.0.0.0:8086
```

```bash
export BIGTABLE_EMULATOR_HOST=localhost:8086
```

With `BIGTABLE_EMULATOR_HOST` set, the SDK dials the emulator itself: no
credentials and no plugin dial options. Tables and column families are
created with `bigtable.NewAdminClient` or the `cbt` CLI.

In Go tests, [`bttest`](https://pkg.go.dev/cloud.google.com/go/bigtable/bttest)
runs the same emulator in process:

```go
srv, _ := bttest.NewServer("localhost:0")
defer srv.Close()
t.Setenv("BIGTABLE_EMULATOR_HOST", srv.Addr)
```

## Configuration Parameters

All keys live under `boost.factory.bigtable` (env `BOOST_FACTORY_BIGTABLE_*`).

| Key | Description | Default |
|---|---|---|
| `project` | GCP project of the Bigtable instance | required |
| `instance` | Bigtable instance ID | required |
| `credentials.file` | service account key file | empty: Application Default Credentials |

## References

- [Bigtable Documentation](https://cloud.google.com/bigtable/docs)
- [Go client](https://pkg.go.dev/cloud.google.com/go/bigtable)
- [Bigtable emulator](https://cloud.google.com/bigtable/docs/emulator)
- [Boost Framework](https://github.com/xgodev/boost)
