# Boost

[![Go](https://github.com/xgodev/boost/actions/workflows/go.yml/badge.svg)](https://github.com/xgodev/boost/actions/workflows/go.yml)

## Overview

Boost is a modular and extensible framework for Go application development, designed to simplify the creation of robust, scalable, and observable services. The framework provides a comprehensive set of components that can be used independently or combined to build complete applications.

## Key Features

- **Modular Architecture**: Independent components that can be used as needed
- **Dependency Injection**: Flexible system for managing dependencies between components
- **Integrated Observability**: Native support for metrics, logging, and tracing
- **Adapters for Popular Libraries**: Consistent wrappers for various libraries in the Go ecosystem
- **Design Patterns**: Implementations of patterns such as Factory, Middleware, and Wrapper
- **CloudEvents Support**: Native integration with the CloudEvents standard
- **Extensibility**: Easy to add support for new libraries and frameworks

## Project Structure

The project is organized into modular packages, each with specific responsibilities:

- **bootstrap**: Components for application initialization
- **examples**: Framework usage examples
- **extra**: Additional functionalities (health checks, middleware, multiserver)
- **factory**: Factory pattern implementations for various components
- **fx**: Dependency injection and lifecycle management
- **model**: Data structures and interfaces definitions
- **utils**: Various utilities
- **wrapper**: Adapters for external libraries (cache, config, log, publisher)

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Specific dependencies vary according to the components used

### Installation

```bash
go get github.com/xgodev/boost
```

### Basic Example

```go
package main

import (
    "context"
    "github.com/xgodev/boost/bootstrap"
    "github.com/xgodev/boost/wrapper/log"
    "github.com/xgodev/boost/wrapper/config"
)

func main() {
    ctx := context.Background()
    
    // Initialize configuration
    cfg := config.New()
    
    // Initialize logger
    logger := log.New(cfg)
    
    // Start application
    app := bootstrap.NewApp(cfg, logger)
    app.Start(ctx)
    
    // Wait for termination signal
    app.WaitForShutdown()
}
```

## Main Components

### Bootstrap

The Bootstrap package provides components for application initialization, including support for serverless functions and event processing.

[Bootstrap Documentation](./bootstrap/README.md)

### Factory

The Factory package implements the Factory design pattern for various components, facilitating the creation and configuration of complex objects.

[Factory Documentation](./factory/README.md)

### Wrapper

The Wrapper package provides adapters for external libraries, including cache, configuration, logging, and message publishing.

[Wrapper Documentation](./wrapper/README.md)

### FX

The FX package provides dependency injection and component lifecycle management functionalities.

[FX Documentation](./fx/README.md)

### Extra

The Extra package includes additional functionalities such as health checks, middleware, and support for multiple servers.

[Extra Documentation](./extra/README.md)

## Integration with Popular Technologies

Boost offers extensive integrations with a wide range of technologies:

### Cloud Providers
- **AWS**: AWS SDK integration for various services
- **Google Cloud Platform**: 
  - Pub/Sub
  - BigQuery
  - Firestore
  - API integration
  - gRPC

### Databases
- **MongoDB**: Native driver integration
- **PostgreSQL**: pgx driver
- **Oracle**: godror driver
- **CockroachDB**: Compatible with PostgreSQL drivers
- **Cassandra**: gocql driver
- **BuntDB**: In-memory key/value database
- **MemDB**: In-memory database with immutable radix tree

### Caching
- **Redis**: go-redis client
- **BigCache**: In-memory caching system
- **FreeCache**: Zero GC cache library

### Messaging & Event Systems
- **NATS**: Lightweight messaging system
- **Kafka**: Confluent Kafka Go client
- **CloudEvents**: SDK for CloudEvents standard
- **Go Cloud Pub/Sub**: Cloud-agnostic pub/sub API

### Web & API
- **Echo**: High performance web framework
- **gRPC**: High performance RPC framework
- **GraphQL**: GraphQL server implementation
- **Swagger**: API documentation with echo-swagger

### Observability
- **Prometheus**: Metrics and monitoring
- **OpenTelemetry**: Observability framework
- **Datadog**: APM and monitoring
- **Zap**: Structured logging
- **Zerolog**: Zero-allocation JSON logger
- **Logrus**: Structured logger

### Configuration
- **Koanf**: Lightweight, extensible configuration management
- **Cobra**: CLI application library
- **Viper**: Configuration solution (via Koanf)

### Resilience & Patterns
- **Hystrix**: Latency and fault tolerance library
- **Ants**: Goroutine pool
- **FTP**: FTP client
- **Kubernetes**: Client-go integration
- **Vault**: HashiCorp Vault client

### Serialization
- **JSON**: Multiple implementations (standard, go-json)
- **MessagePack**: Binary serialization
- **GOB**: Go's native binary format

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository
2. Create a branch for your feature (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the terms included in the [LICENSE](./LICENSE) file.

## Contact

For questions, suggestions, or contributions, please open an issue in the GitHub repository.

---

Developed with ❤️ by the xgodev community.
