# Zerolog Integration for Boost

## Overview

The Zerolog integration in Boost provides a streamlined wrapper around the popular [zerolog](https://github.com/rs/zerolog) structured logging library. This package enables developers to quickly set up structured, leveled logging with multiple output options while seamlessly integrating with Boost's ecosystem.

## Features

- **Structured Logging**: JSON-based structured logging for better parsing and analysis
- **Multiple Output Formats**: Support for text, JSON, and AWS CloudWatch formats
- **Console and File Outputs**: Configurable console and file logging
- **Log Rotation**: Built-in log rotation with size and age limits
- **Compression**: Optional log file compression
- **Configuration Management**: Support for file-based and environment variable configuration
- **Boost Integration**: Seamless integration with Boost's logging interface

## Installation

```go
import (
    "github.com/xgodev/boost/factory/contrib/rs/zerolog/v1"
)
```

## Basic Usage

```go
package main

import (
    "context"
    
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/rs/zerolog/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Initialize Boost
    boost.Start()
    
    // Create a new zerolog logger
    logger := zerolog.NewLogger()
    
    // Set as default logger
    log.SetLogger(logger)
    
    // Use the logger
    log.Info("Hello from zerolog!")
    log.With("key", "value").Warn("Structured logging example")
    
    // Use with context
    ctx := log.NewContext(context.Background(), logger)
    log.FromContext(ctx).Error("Error message with context")
}
```

## Configuration Options

The Zerolog integration can be configured through configuration files or environment variables:

### Environment Variables

```
BOOST_FACTORY_ZEROLOG_LEVEL=INFO
BOOST_FACTORY_ZEROLOG_CONSOLE_ENABLED=true
BOOST_FACTORY_ZEROLOG_FILE_ENABLED=true
BOOST_FACTORY_ZEROLOG_FILE_PATH=/var/log
BOOST_FACTORY_ZEROLOG_FILE_NAME=application.log
BOOST_FACTORY_ZEROLOG_FILE_MAXSIZE=100
BOOST_FACTORY_ZEROLOG_FILE_COMPRESS=true
BOOST_FACTORY_ZEROLOG_FILE_MAXAGE=28
BOOST_FACTORY_ZEROLOG_FORMATTER=JSON
```

### Configuration File (YAML)

```yaml
boost:
  factory:
    zerolog:
      level: INFO
      console:
        enabled: true
      file:
        enabled: true
        path: /var/log
        name: application.log
        maxsize: 100
        compress: true
        maxage: 28
      formatter: JSON
```

### Available Options

| Option | Description | Default |
|--------|-------------|---------|
| `level` | Log level (DEBUG, INFO, WARN, ERROR) | INFO |
| `console.enabled` | Enable console logging | true |
| `file.enabled` | Enable file logging | false |
| `file.path` | Log file directory | /tmp |
| `file.name` | Log file name | application.log |
| `file.maxsize` | Maximum log file size in MB | 100 |
| `file.compress` | Enable log file compression | true |
| `file.maxage` | Maximum log file age in days | 28 |
| `formatter` | Log format (TEXT, JSON, AWS_CLOUD_WATCH) | TEXT |

## Advanced Usage

### Custom Configuration

```go
package main

import (
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/rs/zerolog/v1"
    "github.com/xgodev/boost/wrapper/config"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Initialize Boost with custom config file
    boost.StartWithConfigFile("config.yaml")
    
    // Create a new zerolog logger with options from config
    logger := zerolog.NewLogger()
    
    // Set as default logger
    log.SetLogger(logger)
    
    // Use the logger
    log.Info("Configured zerolog is ready!")
}
```

### Integration with Boost Components

The Zerolog integration works seamlessly with other Boost components:

```go
package main

import (
    "context"
    
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/rs/zerolog/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Initialize Boost
    boost.Start()
    
    // Create a new zerolog logger
    logger := zerolog.NewLogger()
    
    // Create context with logger
    ctx := log.NewContext(context.Background(), logger)
    
    // Use with other Boost components that accept context
    service := myservice.New(ctx)
    
    // The service will use the zerolog logger from context
    service.DoSomething()
}
```

## Best Practices

1. **Use Structured Logging**: Add context to your logs with key-value pairs
   ```go
   log.With("user_id", userId).With("action", "login").Info("User logged in")
   ```

2. **Choose the Right Log Level**: Use appropriate log levels for different types of messages
   ```go
   log.Debug("Detailed debugging information")
   log.Info("Normal application behavior")
   log.Warn("Something unexpected but not critical")
   log.Error("Something failed but application can continue")
   ```

3. **Configure for Production**: In production environments, consider:
   - Using JSON formatter for machine parsing
   - Enabling file logging with rotation
   - Setting appropriate log levels

4. **Pass Logger via Context**: Use context to pass the logger to functions and services
   ```go
   func ProcessRequest(ctx context.Context, req Request) {
       logger := log.FromContext(ctx)
       logger.With("request_id", req.ID).Info("Processing request")
   }
   ```

5. **Log Initialization**: Log when your application starts and important components initialize
   ```go
   log.With("version", appVersion).Info("Application starting")
   ```

## Contributing

Contributions to improve the Zerolog integration are welcome. Please follow the Boost project's contribution guidelines.

## License

This package is part of the Boost project and is subject to its license terms.
