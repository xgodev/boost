# Extra

The Extra package provides additional and complementary functionalities that enrich the Boost framework but don't fit directly into the main packages. This package contains components that add value to the framework, offering solutions for specific use cases.

## Overview

Extra includes implementations for common needs in modern applications, such as health checks, middleware for request interception, and support for multiple servers. These components are designed to be modular and easily integrable with the rest of the framework.

## Main Components

### Health

The `health` subpackage provides implementations for application and service health checks. It includes:

- Endpoints for status verification
- Checkers for different system components
- Aggregators to consolidate results from multiple checks

This component is essential for monitoring and operating applications in production environments, allowing quick detection of problems and facilitating the implementation of recovery strategies.

### Middleware

The `middleware` subpackage offers middleware implementations for intercepting and modifying execution flows. It includes middleware for:

- Request and response logging
- Authentication and authorization
- Error handling
- Metrics and monitoring
- Rate limiting

These middleware can be composed and applied at different points in the application, allowing the implementation of cross-cutting functionalities in a modular and reusable way.

### Multiserver

The `multiserver` subpackage provides support for running multiple servers in a single application. It allows:

- Coordinated initialization and management of servers
- Resource sharing between servers
- Graceful shutdown of all servers

This component is useful for applications that need to expose multiple interfaces, such as REST APIs, gRPC, and administration interfaces, in a single instance.

## How to Use

To use the Extra package in your application:

1. Import the specific subpackage containing the desired functionalities
2. Configure the components as needed
3. Integrate the components with the rest of your application

## Integration with Other Packages

Extra integrates with other components of the Boost framework:

- Bootstrap: For component initialization and configuration
- Factory: For creating instances of specific implementations
- Wrapper: For accessing adapters to external libraries
- Model: For defining interfaces and structures

## Extensibility

The Extra package was designed to be easily extensible. To add new functionalities:

1. Identify the appropriate subpackage or create a new one if necessary
2. Implement the components following existing patterns
3. Add appropriate tests and documentation

## Contribution

Contributions to the Extra package are welcome. When contributing, make sure to follow the code standards and add appropriate tests for new functionalities.
