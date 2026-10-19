# Factory

The Factory package implements the Factory design pattern for various components of the Boost framework, facilitating the creation and configuration of complex objects. This package is fundamental to the modular and extensible architecture of the framework.

## Overview

Factory provides a standardized approach to instantiate objects, allowing client code to request objects without knowing the specific details of how they are created. This promotes decoupling between components and facilitates system extensibility.

## Main Components

### Core

The `core` subpackage contains factory implementations for fundamental system components, organized into categories:

- **Database**: Factories for database connections and operations
- **Net**: Components related to networks and communication
- **Time**: Utilities related to time and scheduling

### Contrib

The `contrib` directory contains factory implementations for integration with third-party libraries and frameworks, including:

- gRPC: For communication between services
- MCP (modelcontextprotocol/go-sdk): Model Context Protocol server and client
- Google Cloud Platform (GCP): Integration with GCP services
- Confluent: For working with Apache Kafka
- Ollama: Integration with AI models
- Spf13/Cobra: For creating command-line interfaces

### Local

The `local` subpackage provides specific implementations for local use or development environments, facilitating testing and rapid prototyping.

## How to Use

To use the Factory package in your application:

1. Import the specific subpackage that contains the desired factory
2. Configure the necessary parameters
3. Use the factory to create instances of the required objects

## Plugins

Factories accept plugins as a variadic argument and run them right after
building the component. The contract lives in this package:

```go
// Plugin customizes a component T built by a factory.
type Plugin[T any] func(context.Context, T) error

// ApplyPlugins runs plugins in order and stops at the first error.
func ApplyPlugins[T any](ctx context.Context, target T, plugins ...Plugin[T]) error
```

The per-factory plugin types are aliases of it, so a plugin written against
either name is interchangeable:

| Factory | Alias |
|---|---|
| `contrib/redis/go-redis/v9` | `Plugin = factory.Plugin[*redis.Client]`, `ClusterPlugin = factory.Plugin[*redis.ClusterClient]` |
| `contrib/labstack/echo/v4`, `v5` | `Plugin = factory.Plugin[*Server]` |
| `contrib/go-resty/resty/v2` | `Plugin = factory.Plugin[*resty.Client]` |
| `contrib/elastic/go-elasticsearch/v8` | `Plugin = factory.Plugin[*elasticsearch.Client]` |
| `contrib/gocql/gocql/v1` | `Plugin = factory.Plugin[*gocql.Session]` |
| `contrib/aws/aws-sdk-go-v2/v1` | `Plugin = factory.Plugin[*aws.Config]` |
| `contrib/go.mongodb.org/mongo-driver/v1`, `v2` | `ClientOptionsPlugin = factory.Plugin[*options.ClientOptions]`, `ClientPlugin = factory.Plugin[*mongo.Client]` |

Factories whose plugins have a different shape (gRPC, NATS, Hystrix, the
mongo `Plugin` pair) keep their own types.

Plugin packages sit under each factory's `plugins/` directory:
`contrib/` integrates a third-party library, `native/` wraps the upstream
library's own middleware, `extra/` is boost-authored and generic, and
`local/` wires another boost package (health, log, cache).

## Integration with Other Packages

The Factory integrates with other components of the Boost framework:

- Bootstrap: For application initialization
- Wrapper: For configuration and logging
- Model: For data structure definitions

## Extensibility

The Factory package was designed to be easily extensible. To add support for new libraries or frameworks:

1. Create a new subpackage in `contrib` with the name of the library
2. Implement the necessary factory interfaces
3. Add appropriate tests and documentation

## Contribution

Contributions to the Factory package are welcome. When contributing, make sure to follow the code standards and add appropriate tests for new functionalities.
