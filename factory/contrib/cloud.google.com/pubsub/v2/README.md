# Google Cloud Pub/Sub Factory (v2)

## Overview

Builds a [`cloud.google.com/go/pubsub/v2`](https://pkg.go.dev/cloud.google.com/go/pubsub/v2)
client from Boost config: API options, gRPC options and OpenTelemetry.
v2 is the default Pub/Sub version of Boost; the v1 factory
(`../v1`) stays for apps that pin `version: v1`.

## Features

- **Client from config**: `NewClient` reads `boost.factory.gcp.pubsub.*`
- **Custom config path or options**: `NewClientWithConfigPath`, `NewClientWithOptions`
- **gRPC plugins**: any `clientgrpc.Plugin` adds dial and call options
- **OpenTelemetry**: gRPC traces and metrics plus Pub/Sub tracing, on by default
- **Logging**: the client creation is logged through `wrapper/log`

## Choosing the version

The fx modules pick the factory by `boost.factory.gcp.pubsub.version`:
`v2` (default) or `v1`. Any other value fails the boot. The publisher driver
and the function adapter have their own version keys, also `v2` by default:

| Key | Selects |
|---|---|
| `boost.factory.gcp.pubsub.version` | the `*pubsub.Client` of `fx/modules/factory/.../pubsub` |
| `boost.wrapper.publisher.driver.pubsub.version` | the `publisher.Driver` of `fx/modules/wrapper/publisher/...` |
| `boost.bootstrap.function.adapter.pubsub.version` | the `@Function` adapter of `fx/modules/bootstrap/function/...` |

To stay on v1, set the three keys to `v1`.

## Usage

### Basic Usage

```go
package main

import (
    "context"

    gpubsub "cloud.google.com/go/pubsub/v2"
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v2"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    boost.Start()

    ctx := context.Background()

    client, err := pubsub.NewClient(ctx)
    if err != nil {
        log.Fatalf("failed to create Pub/Sub client: %v", err)
    }
    defer client.Close()

    publisher := client.Publisher("my-topic")
    defer publisher.Stop()

    id, err := publisher.Publish(ctx, &gpubsub.Message{
        Data: []byte("Hello, Pub/Sub!"),
    }).Get(ctx)
    if err != nil {
        log.Fatalf("failed to publish message: %v", err)
    }

    log.Infof("published message with ID: %s", id)
}
```

### Custom Configuration Path

```go
client, err := pubsub.NewClientWithConfigPath(ctx, "myapp.pubsub")
```

The path must have the same keys as the default root; `pubsub.ConfigAdd(path)`
registers them with their defaults.

### Custom Options

```go
options, err := pubsub.NewOptions()
if err != nil {
    log.Fatalf("failed to load options: %v", err)
}

options.APIOptions.ProjectID = "my-gcp-project"
options.Otel.Enabled = false

client, err := pubsub.NewClientWithOptions(ctx, options)
```

### Using Plugins

```go
import (
    clientgrpc "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
    "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client/plugins/native/compressor"
)

client, err := pubsub.NewClient(ctx, compressor.Register)
```

A plugin is a `clientgrpc.Plugin`:
`func(ctx context.Context) ([]grpc.DialOption, []grpc.CallOption)`.

## Configuration Parameters

All keys live under `boost.factory.gcp.pubsub`.

| Key | Description | Default |
|---|---|---|
| `version` | Pub/Sub version the fx modules use: `v1` or `v2` | `v2` |
| `otel.enabled` | gRPC OpenTelemetry traces/metrics and Pub/Sub tracing | `true` |
| `apiOptions.*` | project, credentials, endpoint, scopes ([`api/v0`](../../api/v0)) | application default credentials |
| `grpcOptions.*` | message sizes, keepalive, windows ([`grpc/v1`](../../grpc/v1)) | gRPC defaults |

## Example: Publishing and Subscribing

v2 splits data and admin calls: `Publisher`/`Subscriber` send and receive,
`TopicAdminClient`/`SubscriptionAdminClient` create and delete resources.

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"

    gpubsub "cloud.google.com/go/pubsub/v2"
    "cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v2"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    boost.Start()

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    client, err := pubsub.NewClient(ctx)
    if err != nil {
        log.Fatalf("failed to create client: %v", err)
    }
    defer client.Close()

    topic := fmt.Sprintf("projects/%s/topics/my-topic", client.Project())
    if _, err := client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{Name: topic}); err != nil {
        log.Fatalf("failed to create topic: %v", err)
    }

    sub := fmt.Sprintf("projects/%s/subscriptions/my-subscription", client.Project())
    if _, err := client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
        Name:               sub,
        Topic:              topic,
        AckDeadlineSeconds: 20,
    }); err != nil {
        log.Fatalf("failed to create subscription: %v", err)
    }

    publisher := client.Publisher(topic)
    defer publisher.Stop()

    id, err := publisher.Publish(ctx, &gpubsub.Message{
        Data:       []byte("Hello, Pub/Sub!"),
        Attributes: map[string]string{"origin": "boost-example"},
    }).Get(ctx)
    if err != nil {
        log.Fatalf("failed to publish message: %v", err)
    }
    log.Infof("published message with ID: %s", id)

    rctx, stop := context.WithCancel(ctx)
    err = client.Subscriber(sub).Receive(rctx, func(_ context.Context, msg *gpubsub.Message) {
        log.Infof("received message: %s %v", msg.Data, msg.Attributes)
        msg.Ack()
        stop()
    })
    if err != nil && !errors.Is(err, context.Canceled) {
        log.Fatalf("failed to receive messages: %v", err)
    }
}
```

## Best Practices

1. **Project ID**: set `apiOptions.projectID` in production
2. **Authentication**: service account with the minimum Pub/Sub roles
3. **Stop publishers**: `Publisher.Stop()` flushes pending messages
4. **Ack deadline**: match it to the processing time of a message
5. **Admin calls**: create topics and subscriptions in provisioning, not on every boot

## References

- [Google Cloud Pub/Sub Documentation](https://cloud.google.com/pubsub/docs)
- [Go client v2](https://pkg.go.dev/cloud.google.com/go/pubsub/v2)
- [Migrating from v1 to v2](https://github.com/googleapis/google-cloud-go/blob/main/pubsub/MIGRATING.md)
- [Boost Framework](https://github.com/xgodev/boost)
