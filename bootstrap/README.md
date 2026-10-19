# Bootstrap

The Bootstrap package is responsible for providing components and utilities for initializing applications in the Boost framework. This package serves as a foundation for configuring and starting applications in a standardized and efficient manner.

## Overview

Bootstrap acts as an entry point for applications built with the Boost framework, offering structures and patterns for service initialization. It defines fundamental constants and provides mechanisms for initial application configuration.

## Main Components

### Function

The `function` subpackage provides structures and utilities for working with functions as processing units, especially in the context of cloud events (Cloud Events). It includes:

- Handlers for event processing
- Middleware for intercepting and modifying execution flows
- Adapters for different protocols and formats
- Wrappers for functions that facilitate integration with the Boost ecosystem

The cloudevents nats adapter has no responder: the event a function returns
goes to its middlewares, like `publisher`, and never back on the subject.

### Examples

`function/examples` holds one builder app per adapter and middleware
(cloudevents http and nats, nats, pubsub, middlewares, recovery): an annotated
`@Function` handler, the committed `gen/`, the generated `cmd/worker/main.go`
and a test that runs the worker in process.

## How to Use

To use the Bootstrap package in your application:

1. Import the package in your code
2. Configure the necessary parameters
3. Use the Function components to process events or start services

## Integration with Other Packages

Bootstrap integrates with other components of the Boost framework, such as:

- Wrapper: for configuration and logging
- Factory: for component creation
- Extra: for additional functionalities like middleware and health checks

## Contribution

Contributions to the Bootstrap package are welcome. When contributing, make sure to follow the code standards and add appropriate tests for new functionalities.
