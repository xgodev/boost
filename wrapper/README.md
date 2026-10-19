# Wrapper

The Wrapper package provides adapters and abstractions for external libraries, allowing standardized and flexible integration with the Boost framework. This package is essential for ensuring interoperability and consistency when using third-party libraries.

## Overview

Wrapper acts as an abstraction layer between application code and external libraries, offering consistent interfaces and facilitating the replacement of underlying implementations without affecting client code. This promotes code maintainability and testability.

## Main Components

### Cache

The `cache` subpackage provides abstractions for caching systems, including:

- **Codec**: Serializers and deserializers for different formats (binary, gob, json, string)
- **Driver**: Implementations for different cache providers (Redis, BigCache, FreeCache)
- **Plugins**: Extensions for additional cache functionalities

This component allows efficient storage and retrieval of cached data, with support for multiple formats and backends.

### Config

The `config` subpackage offers a unified interface for application configuration, with support for:

- Multiple configuration sources (files, environment variables, flags)
- Different formats (YAML, JSON, ENV)
- Value validation and transformation

This facilitates configuration management across different environments and deployment scenarios.

`config.Add` keeps one entry per key: a key added again, like the keys the v1
and v2 pubsub adapters share, replaces its entry, so the flag parsers never see
a key twice. Environment variables map to keys by `_` → `.`, with `__` joining
the words of a camelCase key: `BOOST_FACTORY_GCP_PUBSUB_API__OPTIONS_EMULATOR__HOST`
sets `boost.factory.gcp.pubsub.apiOptions.emulatorHost`.

### Log

The `log` subpackage provides abstractions for logging systems, with support for:

- Multiple log levels (debug, info, warn, error)
- Different output formats (text, JSON, CloudWatch)
- Integration with popular libraries (Zap, Zerolog, Logrus)

This component allows consistent and configurable recording of application events and information.

### Publisher

The `publisher` subpackage offers abstractions for message publishing/subscription systems, including:

- Drivers for different providers (Google Cloud Pub/Sub, Kafka, NATS)
- Middleware for message interception and modification
- Utilities for monitoring and metrics

This facilitates asynchronous communication between components and services.

`Publisher.Publish` returns an error when the driver fails or when any event
fails. A per-event failure comes back as one `Internal` error that joins every
failed event's error (tagged with its event ID), so `errors.Is`/`errors.As`
still reach each cause. Middlewares keep receiving the per-event
`[]PublishOutput` (with each event's `Error`) and a nil chain error; the
aggregate error is built after the middleware chain.

The Kafka driver (`driver/contrib/confluentinc/confluent-kafka-go/v2`) returns
only after every produced message has its delivery report, bounded by the
`ctx` passed to `Publish`: a failed delivery (`TopicPartition.Error`, e.g.
message timeout or unknown partition) or a `ctx` that ends before the report
arrives becomes that event's `Error`. Reports go to a delivery channel per
`Publish` call, so they no longer reach the producer's `Events()` channel; the
driver logs each delivery at its own `log.level`.

## How to Use

To use the Wrapper package in your application:

1. Import the specific subpackage containing the desired abstraction
2. Configure the specific implementation you want to use
3. Use the abstract interface in your code, independent of the underlying implementation

## Integration with Other Packages

Wrapper integrates with other components of the Boost framework:

- Bootstrap: For component initialization and configuration
- Factory: For creating instances of specific implementations
- Extra: For additional functionalities like middleware

## Extensibility

The Wrapper package was designed to be easily extensible. To add support for new libraries:

1. Create a new subpackage in `contrib` with the name of the library
2. Implement the necessary abstract interfaces
3. Add appropriate tests and documentation

## Contribution

Contributions to the Wrapper package are welcome. When contributing, make sure to follow the code standards and add appropriate tests for new functionalities.
