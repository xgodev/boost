# FX

The FX package provides dependency injection and component lifecycle management functionalities for applications built with the Boost framework. This package is essential for building modular and testable applications.

## Overview

FX simplifies the creation and management of dependencies between components, allowing developers to define reusable modules and compose complex applications in an organized manner. It draws inspiration from dependency injection frameworks, adapted for the Go ecosystem.

## Main Components

### Modules

The `modules` subpackage contains implementations of pre-defined modules that can be used to compose applications. These modules encapsulate common functionalities and can be easily integrated into different applications.

The available modules include integrations with various components of the Boost ecosystem, such as caching systems, configuration, logging, and message publishing. Each module is designed to be independent and easily combinable with other modules.

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
