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
