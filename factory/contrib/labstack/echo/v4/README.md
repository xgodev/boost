# Echo Integration for Boost

## Overview

The Echo integration in Boost provides a robust wrapper around the popular [Echo web framework](https://echo.labstack.com/) (v4), offering simplified server initialization, configuration management, middleware integration, and plugin support. This package enables developers to quickly set up HTTP, H2C, and TLS servers with a clean API while leveraging Boost's ecosystem.

## Features

- **Simplified Server Creation**: Easy initialization with sensible defaults
- **Configuration Management**: Support for file-based configuration
- **Plugin Architecture**: Extensible design with plugin support
- **Middleware Integration**: Built-in support for Echo middleware
- **Protocol Support**: HTTP, H2C, and TLS (including auto-certificate generation)
- **Routing API**: Complete access to Echo's routing capabilities
- **Boost Integration**: Seamless integration with Boost logging and other components

## Installation

```go
import (
    "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
)
```

## Basic Usage

### Creating a Server

```go
// Create a new server with default options
server, err := echo.NewServer(ctx)
if err != nil {
    log.Fatal(err)
}

// Define routes
server.GET("/hello", func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
})

// Start the server
server.Serve(ctx)
```

### Using Configuration File

```go
// Create a server with options from config file
server, err := echo.NewServerWithConfigPath(ctx, "/path/to/config.yaml")
if err != nil {
    log.Fatal(err)
}

// Define routes and start server
// ...
```

### Using Plugins

```go
// Create plugins
corsPlugin := echoplugins.CORS(ctx)
recoverPlugin := echoplugins.Recover(ctx)

// Create server with plugins
server, err := echo.NewServer(ctx, corsPlugin, recoverPlugin)
if err != nil {
    log.Fatal(err)
}

// Define routes and start server
// ...
```

## Configuration Options

The Echo server can be configured with various options:

```go
options := &echo.Options{
    Port:         8080,
    Protocol:     "HTTP", // "HTTP", "H2C", or "TLS"
    HideBanner:   true,
    DisableHTTP2: false,
    TLS: echo.TLSOptions{
        Enabled: false,
        Type:    "FILE", // "FILE" or "AUTO"
        File: echo.TLSFileOptions{
            Cert: "/path/to/cert.pem",
            Key:  "/path/to/key.pem",
        },
        Auto: echo.TLSAutoOptions{
            Host: "example.com",
        },
    },
}

server := echo.NewServerWithOptions(ctx, options)
```

## Available Plugins

The Echo integration includes several built-in plugins organized into categories:

### Native Plugins

These plugins wrap Echo's native middleware:

- **BodyDump**: Dumps request and response bodies
- **BodyLimit**: Sets the maximum allowed request body size
- **CORS**: Configures Cross-Origin Resource Sharing
- **GZIP**: Compresses responses using gzip
- **Recover**: Recovers from panics and returns 500 error
- **RequestID**: Adds a unique request ID to each request

### Extra Plugins

Additional functionality beyond Echo's native capabilities:

- **ErrorHandler**: Custom error handling
- **Semaphore**: Request rate limiting

### Local Plugins

Boost-specific integrations:

- **Health**: Health check endpoints
- **MultiServer**: Run multiple Echo servers
- **RestResponse**: Standardized REST response format
- **Log**: Integration with Boost logging

### Contrib Plugins

Third-party integrations:

- **Echo-Pprof**: Profiling endpoints
- **Prometheus**: Metrics collection
- **Swagger**: API documentation

## Example: Complete Server with Middleware

```go
package main

import (
    "context"
    "net/http"

    "github.com/labstack/echo/v4"
    echoserver "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
    "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/cors"
    "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/recover"
    "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
)

func main() {
    ctx := context.Background()
    
    // Create plugins
    corsPlugin := cors.New(ctx)
    recoverPlugin := recover.New(ctx)
    logPlugin := log.New(ctx)
    
    // Create server with plugins
    server, err := echoserver.NewServer(ctx, corsPlugin, recoverPlugin, logPlugin)
    if err != nil {
        panic(err)
    }
    
    // Define routes
    server.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    
    // Group routes
    api := server.Group("/api")
    api.GET("/users", getUsers)
    api.POST("/users", createUser)
    
    // Start server
    server.Serve(ctx)
}

func getUsers(c echo.Context) error {
    return c.JSON(http.StatusOK, []string{"user1", "user2"})
}

func createUser(c echo.Context) error {
    // Implementation...
    return c.JSON(http.StatusCreated, map[string]string{"status": "created"})
}
```

## Integration with Boost

The Echo integration works seamlessly with other Boost components:

- **Logging**: Automatically integrates with Boost's logging system
- **Configuration**: Uses Boost's configuration management
- **Metrics**: Can be integrated with Prometheus for metrics collection
- **Health Checks**: Built-in health check endpoints

## Best Practices

1. **Use Plugins for Middleware**: Organize middleware as plugins for better code organization
2. **Group Related Routes**: Use the Group method to organize related endpoints
3. **Leverage Configuration Files**: Store server configuration in external files
4. **Implement Proper Error Handling**: Use the ErrorHandler plugin for consistent error responses
5. **Add Request Logging**: Use the Log plugin to log all requests
6. **Set Appropriate Timeouts**: Configure timeouts to prevent resource exhaustion
7. **Use TLS in Production**: Enable TLS for production environments

## Contributing

Contributions to improve the Echo integration are welcome. Please follow the Boost project's contribution guidelines.

## License

This package is part of the Boost project and is subject to its license terms.
