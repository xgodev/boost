# Hystrix Factory

## Overview

The Hystrix Factory provides integration with the [Hystrix-Go](https://github.com/afex/hystrix-go) library, implementing the circuit breaker pattern for Go applications within the Boost framework. This factory enables resilient and fault-tolerant systems by preventing cascading failures in distributed environments.

## Features

- **Command Configuration**: Easy configuration of Hystrix commands through Boost's configuration system
- **Circuit Breaker Pattern**: Implementation of the circuit breaker pattern to handle failures gracefully
- **Configurable Parameters**: Fine-grained control over timeout, concurrency, error thresholds, and recovery windows
- **Logging Integration**: Seamless integration with Boost's logging system

## Usage

### Basic Configuration

```go
package main

import (
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/afex/hystrix-go/v0"
)

func main() {
    // Initialize Boost
    boost.Start()
    
    // Configure a Hystrix command
    hystrix.ConfigureCommand("my_command")
    
    // Use the command in your application
    // ...
}
```

### Configuring Multiple Commands

```go
// Configure multiple commands at once
commands := []string{"database_query", "api_request", "authentication"}
hystrix.ConfigureCommands(commands)
```

### Custom Configuration

You can customize Hystrix command parameters through Boost's configuration system:

```
# In your configuration file or environment variables
boost.factory.hystrix.commands.my_command.timeout=5000
boost.factory.hystrix.commands.my_command.maxConcurrentRequests=10
boost.factory.hystrix.commands.my_command.requestVolumeThreshold=20
boost.factory.hystrix.commands.my_command.errorPercentThreshold=25
boost.factory.hystrix.commands.my_command.sleepWindow=10000
```

## Configuration Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `timeout` | How long to wait for command to complete (milliseconds) | 10000 |
| `maxConcurrentRequests` | Maximum number of concurrent requests | 20 |
| `requestVolumeThreshold` | Minimum requests needed before circuit can trip | 10 |
| `errorPercentThreshold` | Error percentage to open the circuit | 5 |
| `sleepWindow` | Time to wait after circuit opens before testing recovery (milliseconds) | 5000 |

## Integration with Other Boost Components

The Hystrix Factory integrates with:

- **Config Wrapper**: For loading and managing configuration
- **Log Wrapper**: For logging circuit breaker events and errors

## Plugin System

The factory supports plugins for extending functionality. Custom plugins can be added to the `plugins` directory and registered with the factory.

## Example

```go
package main

import (
    "errors"
    "github.com/afex/hystrix-go/hystrix"
    hystrixfactory "github.com/xgodev/boost/factory/contrib/afex/hystrix-go/v0"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Configure the command
    hystrixfactory.ConfigureCommand("get_user_profile")
    
    // Use the command
    err := hystrix.Do("get_user_profile", func() error {
        // This function will be executed if the circuit is closed
        return callExternalService()
    }, func(err error) error {
        // This function will be executed if the circuit is open
        log.Errorf("Circuit open, using fallback: %v", err)
        return getFallbackUserProfile()
    })
    
    if err != nil {
        log.Errorf("Error: %v", err)
    }
}

func callExternalService() error {
    // Actual implementation to call external service
    return nil
}

func getFallbackUserProfile() error {
    // Fallback implementation
    return errors.New("using cached profile")
}
```

## Best Practices

1. **Command Naming**: Use descriptive names for commands that reflect their purpose
2. **Timeout Configuration**: Set appropriate timeouts based on expected response times
3. **Error Thresholds**: Adjust error thresholds based on your application's tolerance for failures
4. **Monitoring**: Monitor circuit breaker states to understand system health

## References

- [Hystrix-Go GitHub Repository](https://github.com/afex/hystrix-go)
- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html)
