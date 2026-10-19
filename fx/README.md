# FX

The FX package provides dependency injection and component lifecycle management functionalities for applications built with the Boost framework. This package is essential for building modular and testable applications.

## Overview

FX simplifies the creation and management of dependencies between components, allowing developers to define reusable modules and compose complex applications in an organized manner. It draws inspiration from dependency injection frameworks, adapted for the Go ecosystem.

## Main Components

### Modules

The `modules` subpackage contains implementations of pre-defined modules that can be used to compose applications. These modules encapsulate common functionalities and can be easily integrated into different applications.

The available modules include integrations with various components of the Boost ecosystem, such as caching systems, configuration, logging, and message publishing. Each module is designed to be independent and easily combinable with other modules.

### Factory modules

Each `fx/modules/factory/contrib/<path>` mirrors `factory/contrib/<path>` and wraps its default constructor. Plugins are optional (`[]Plugin` in the graph). Resources with a close/shutdown are released on app stop.

| Module | Entry points | Provides | Requires from the app |
|---|---|---|---|
| `allegro/bigcache/v3` | `Module` | `*bigcache.BigCache` | — |
| `aws/aws-sdk-go-v2/v1` | `Module` | `aws.Config` | — |
| `cloud.google.com/bigquery/v1` | `Module` | `*bigquery.Client` | — |
| `cloud.google.com/firestore/v1` | `Module` | `*firestore.Client` | — |
| `cloud.google.com/pubsub/v1`, `v2` | `Module` | `*pubsub.Client` | — |
| `confluentinc/confluent-kafka-go/v2` | `ProducerModule`, `ConsumerModule` | kafka producer / consumer | — |
| `coocood/freecache/v1` | `Module` | `*freecache.Cache` | — |
| `datadog/dd-trace-go/v1` | `TracerModule`, `ProfilerModule` | starts tracer/profiler | — |
| `elastic/go-elasticsearch/v8`, `v9` | `Module`, `BulkModule` | `*elasticsearch.Client`, `esutil.BulkIndexer` | — |
| `go-resty/resty/v2` | `Module` | `*resty.Client` | — |
| `go.mongodb.org/mongo-driver/v1`, `v2` | `Module` | `*mongo.Conn`, `*mongo.Database` | — |
| `go.opentelemetry.io/otel/v1` | `Module` | starts tracer and meter providers | — |
| `gocloud.dev/pubsub/v0` | `TopicModule`, `SubscriptionModule` | `*pubsub.Topic`, `*pubsub.Subscription` | — |
| `gocql/gocql/v1` | `Module` | `*gocql.Session` | — |
| `godror/godror/v0` | `Module` | `*sql.DB` | — |
| `google.golang.org/grpc/v1/client` | `Module` | `*grpc.ClientConn` | — |
| `google.golang.org/grpc/v1/server` | `Module` | `*server.Server`, `*grpc.Server` | — |
| `graphql-go/graphql/v0` | `Module` | `*handler.Handler` | `*graphql.Schema` |
| `hashicorp/memdb/v1` | `Module` | `*memdb.MemDB` | `*memdb.DBSchema` |
| `jackc/pgx/v5` | `Module` | `*sql.DB` | — |
| `jlaffaye/ftp/v0` | `Module` | `*ftp.ServerConn` | — |
| `k8s.io/client-go/v0` | `Module` | `*kubernetes.Clientset` | — |
| `labstack/echo/v4` | `Module` | `*echo.Server`, `*echo.Echo` | — |
| `lovoo/goka/v1` | `Module` | goka emitter | — |
| `mittwald/vaultgo/v0` | `Module` | `*vault.Client` | — |
| `nats-io/nats.go/v1` | `Module` | nats connection | — |
| `panjf2000/ants/v2` | `Module` | `*ants.Wrapper` | — |
| `prometheus/client_golang/v1` | `PusherModule` | starts push gateway pusher (if enabled) | — |
| `redis/go-redis/v9` | `ClientModule`, `ClusterModule` | `*redis.Client`, `*redis.ClusterClient` | — |
| `tidwall/buntdb/v1` | `Module` | `*buntdb.DB` | — |

`godror` and `pgx` both provide `*sql.DB`: compose only one per app.

Pub/Sub v1 and v2 are both composable: `boost.factory.gcp.pubsub.version`
(`v2` by default) turns on one of them. The publisher driver and the function
adapter pick their version the same way, through
`boost.wrapper.publisher.driver.pubsub.version` and
`boost.bootstrap.function.adapter.pubsub.version`. To stay on v1, set the
three keys to `v1`; any value other than `v1` or `v2` fails the boot.

## How to Use

To use the FX package in your application:

1. Import the package and the specific modules you want to use
2. Define your own modules to encapsulate business logic
3. Compose the modules to create your application
4. Use the dependency injection system to access components in your code

FX automatically manages the lifecycle of components, ensuring they are initialized in the correct order and properly finalized when the application is terminated.

## Integration with Other Packages

FX deeply integrates with other components of the Boost framework:

- Bootstrap: For application initialization
- Factory: For creating component instances
- Wrapper: For accessing adapters to external libraries
- Model: For defining interfaces and structures

## Extensibility

The FX package was designed to be easily extensible. To add new modules:

1. Create a new subpackage in `modules` with the name of the module
2. Implement the necessary interfaces to define the module
3. Add appropriate tests and documentation

## Contribution

Contributions to the FX package are welcome. When contributing, make sure to follow the code standards and add appropriate tests for new functionalities.
