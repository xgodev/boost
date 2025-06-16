# Logrus Integration for Boost

## Overview

The Logrus integration in Boost provides a comprehensive wrapper around the popular [Logrus](https://github.com/sirupsen/logrus) structured logging library. This package enables developers to quickly set up structured, leveled logging with multiple output formats and destinations while seamlessly integrating with Boost's ecosystem.

## Features

- **Structured Logging**: Field-based structured logging for better parsing and analysis
- **Multiple Formatters**: Support for Text, JSON, and AWS CloudWatch formats
- **Console and File Outputs**: Configurable console and file logging
- **Log Rotation**: Built-in log file rotation with size and age limits
- **Compression**: Optional log file compression
- **Custom Hooks**: Support for adding custom Logrus hooks
- **Configuration Management**: Support for file-based and environment variable configuration
- **Boost Integration**: Seamless integration with Boost's logging interface

## Installation

```go
import (
    "github.com/xgodev/boost/factory/contrib/sirupsen/logrus/v1"
)
```

## Basic Usage

```go
package main

import (
    "context"
    
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/sirupsen/logrus/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Initialize Boost
    boost.Start()
    
    // Create a new logrus logger
    logger := logrus.NewLogger()
    
    // Set as default logger
    log.SetLogger(logger)
    
    // Use the logger
    log.Info("Hello from logrus!")
    log.With("key", "value").Warn("Structured logging example")
    
    // Use with context
    ctx := log.NewContext(context.Background(), logger)
    log.FromContext(ctx).Error("Error message with context")
}
```

## Configuration Options

The Logrus integration can be configured through configuration files or environment variables:

### Environment Variables

```
BOOST_FACTORY_LOGRUS_CONSOLE_ENABLED=true
BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL=INFO
BOOST_FACTORY_LOGRUS_FILE_ENABLED=true
BOOST_FACTORY_LOGRUS_FILE_LEVEL=INFO
BOOST_FACTORY_LOGRUS_FILE_PATH=/var/log
BOOST_FACTORY_LOGRUS_FILE_NAME=application.log
BOOST_FACTORY_LOGRUS_FILE_MAXSIZE=100
BOOST_FACTORY_LOGRUS_FILE_COMPRESS=true
BOOST_FACTORY_LOGRUS_FILE_MAXAGE=28
BOOST_FACTORY_LOGRUS_TIME_FORMAT="2006/01/02 15:04:05.000"
BOOST_FACTORY_LOGRUS_FORMATTERTYPE=TEXT
```

### Configuration File (YAML)

```yaml
boost:
  factory:
    logrus:
      console:
        enabled: true
        level: INFO
      file:
        enabled: true
        level: INFO
        path: /var/log
        name: application.log
        maxsize: 100
        compress: true
        maxage: 28
      time:
        format: "2006/01/02 15:04:05.000"
      formatterType: TEXT
```

### Available Options

| Option | Description | Default |
|--------|-------------|---------|
| `console.enabled` | Enable console logging | true |
| `console.level` | Console log level (DEBUG, INFO, WARN, ERROR) | INFO |
| `file.enabled` | Enable file logging | false |
| `file.level` | File log level (DEBUG, INFO, WARN, ERROR) | INFO |
| `file.path` | Log file directory | /tmp |
| `file.name` | Log file name | application.log |
| `file.maxsize` | Maximum log file size in MB | 100 |
| `file.compress` | Enable log file compression | true |
| `file.maxage` | Maximum log file age in days | 28 |
| `time.format` | Time format for log entries | 2006/01/02 15:04:05.000 |
| `formatterType` | Log format (TEXT, JSON, CLOUDWATCH) | TEXT |

## Formatters

The Logrus integration includes three built-in formatters:

### Text Formatter

The default formatter that outputs logs in a human-readable format:

```
2025/05/22 14:34:35.123 INFO Hello from logrus!
2025/05/22 14:34:35.124 WARN Structured logging example key=value
```

### JSON Formatter

Outputs logs in JSON format for machine parsing:

```json
{"level":"info","msg":"Hello from logrus!","time":"2025-05-22T14:34:35.123Z"}
{"key":"value","level":"warn","msg":"Structured logging example","time":"2025-05-22T14:34:35.124Z"}
```

### CloudWatch Formatter

Optimized for AWS CloudWatch Logs:

```json
{"level":"info","message":"Hello from logrus!","timestamp":"2025-05-22T14:34:35.123Z"}
{"key":"value","level":"warn","message":"Structured logging example","timestamp":"2025-05-22T14:34:35.124Z"}
```

## Advanced Usage

### Custom Hooks

You can add custom Logrus hooks to extend functionality:

```go
package main

import (
    "github.com/sirupsen/logrus"
    logrusFactory "github.com/xgodev/boost/factory/contrib/sirupsen/logrus/v1"
    "github.com/xgodev/boost/wrapper/log"
)

// Custom hook for Logrus
type CustomHook struct{}

func (h *CustomHook) Levels() []logrus.Level {
    return []logrus.Level{
        logrus.InfoLevel,
        logrus.WarnLevel,
        logrus.ErrorLevel,
    }
}

func (h *CustomHook) Fire(entry *logrus.Entry) error {
    // Custom processing of log entries
    return nil
}

func main() {
    // Create custom hook
    customHook := &CustomHook{}
    
    // Create logger with custom hook
    logger := logrusFactory.NewLogger(customHook)
    
    // Set as default logger
    log.SetLogger(logger)
    
    // Use the logger
    log.Info("This log entry will be processed by the custom hook")
}
```

### Custom Configuration

```go
package main

import (
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/sirupsen/logrus/v1"
    "github.com/xgodev/boost/wrapper/config"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Initialize Boost with custom config file
    boost.StartWithConfigFile("config.yaml")
    
    // Create a new logrus logger with options from config
    logger := logrus.NewLogger()
    
    // Set as default logger
    log.SetLogger(logger)
    
    // Use the logger
    log.Info("Configured logrus is ready!")
}
```

### Integration with Boost Components

The Logrus integration works seamlessly with other Boost components:

```go
package main

import (
    "context"
    
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/sirupsen/logrus/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Initialize Boost
    boost.Start()
    
    // Create a new logrus logger
    logger := logrus.NewLogger()
    
    // Create context with logger
    ctx := log.NewContext(context.Background(), logger)
    
    // Use with other Boost components that accept context
    service := myservice.New(ctx)
    
    // The service will use the logrus logger from context
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

5. **Add Custom Hooks**: Use hooks for specialized logging needs like sending alerts or metrics

## Contributing

Contributions to improve the Logrus integration are welcome. Please follow the Boost project's contribution guidelines.

## License

This package is part of the Boost project and is subject to its license terms.
