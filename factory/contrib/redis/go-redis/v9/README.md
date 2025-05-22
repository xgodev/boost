# Redis Integration for Boost

## Overview

The Redis integration in Boost provides a comprehensive wrapper around the popular [go-redis](https://github.com/redis/go-redis) library (v9), offering simplified client initialization, configuration management, and plugin support. This package enables developers to quickly set up Redis connections with various deployment models (standalone, sentinel, cluster) while leveraging Boost's ecosystem.

## Features

- **Multiple Deployment Models**: Support for standalone Redis, Redis Sentinel, and Redis Cluster
- **Simplified Client Creation**: Easy initialization with sensible defaults
- **Configuration Management**: Support for file-based configuration
- **Plugin Architecture**: Extensible design with plugin support
- **Observability Integration**: Built-in support for Prometheus metrics, Datadog tracing, and OpenTelemetry
- **Health Checks**: Ready-to-use health check implementations
- **Boost Integration**: Seamless integration with Boost logging and other components

## Installation

```go
import (
    "github.com/xgodev/boost/factory/contrib/redis/go-redis/v9"
)
```

## Basic Usage

### Standalone Client

```go
// Create a new Redis client with default options
client, err := redis.NewClient(ctx)
if err != nil {
    log.Fatal(err)
}

// Use the client
val, err := client.Get(ctx, "key").Result()
if err != nil {
    log.Error(err)
}
fmt.Println("value:", val)
```

### Cluster Client

```go
// Create a new Redis Cluster client with default options
cluster, err := redis.NewClusterClient(ctx)
if err != nil {
    log.Fatal(err)
}

// Use the cluster client
val, err := cluster.Get(ctx, "key").Result()
if err != nil {
    log.Error(err)
}
fmt.Println("value:", val)
```

### Using Configuration File

```go
// Create a client with options from config file
client, err := redis.NewClientWithConfigPath(ctx, "/path/to/config.yaml")
if err != nil {
    log.Fatal(err)
}

// Create a cluster client with options from config file
cluster, err := redis.NewClusterClientWithConfigPath(ctx, "/path/to/config.yaml")
if err != nil {
    log.Fatal(err)
}
```

### Using Plugins

```go
// Create health check plugin
healthIntegrator, err := health.NewClientHealth()
if err != nil {
    log.Fatal(err)
}

// Create Prometheus metrics plugin
prometheusPlugin, err := prometheus.NewClientPlugin(ctx)
if err != nil {
    log.Fatal(err)
}

// Create client with plugins
client, err := redis.NewClient(ctx, 
    healthIntegrator.Register,
    prometheusPlugin.Register,
)
if err != nil {
    log.Fatal(err)
}
```

## Configuration Options

The Redis integration can be configured with various options:

```go
options := &redis.Options{
    // Common options
    Password:        "password",
    MaxRetries:      3,
    MinRetryBackoff: time.Second,
    MaxRetryBackoff: time.Second * 3,
    DialTimeout:     time.Second * 5,
    ReadTimeout:     time.Second * 3,
    WriteTimeout:    time.Second * 3,
    PoolSize:        10,
    MinIdleConns:    2,
    PoolTimeout:     time.Second * 4,
    
    // Standalone client options
    Client: redis.ClientOptions{
        Addr:    "localhost:6379",
        Network: "tcp",
        DB:      0,
    },
    
    // Sentinel options
    Sentinel: redis.SentinelOptions{
        MasterName: "mymaster",
        Addrs:      []string{"sentinel1:26379", "sentinel2:26379"},
        Password:   "sentinel-password",
    },
    
    // Cluster options
    Cluster: redis.ClusterOptions{
        Addrs:          []string{"node1:6379", "node2:6379", "node3:6379"},
        MaxRedirects:   3,
        ReadOnly:       false,
        RouteByLatency: false,
        RouteRandomly:  false,
    },
}

// Create client with custom options
client := redis.NewClientWithOptions(ctx, options)
```

## Available Plugins

The Redis integration includes several built-in plugins:

### Health Check Plugins

```go
// For standalone client
healthIntegrator, err := health.NewClientHealth()
if err != nil {
    log.Fatal(err)
}
client, err := redis.NewClient(ctx, healthIntegrator.Register)

// For cluster client
clusterHealthIntegrator, err := health.NewClusterHealth()
if err != nil {
    log.Fatal(err)
}
cluster, err := redis.NewClusterClient(ctx, clusterHealthIntegrator.Register)

// Check health status
healthStatus := h.CheckAll(ctx)
```

### Prometheus Metrics Plugin

```go
// For standalone client
prometheusPlugin, err := prometheus.NewClientPlugin(ctx)
if err != nil {
    log.Fatal(err)
}
client, err := redis.NewClient(ctx, prometheusPlugin.Register)

// For cluster client
prometheusClusterPlugin, err := prometheus.NewClusterPlugin(ctx)
if err != nil {
    log.Fatal(err)
}
cluster, err := redis.NewClusterClient(ctx, prometheusClusterPlugin.Register)
```

### Datadog Tracing Plugin

```go
// For standalone client
datadogPlugin, err := datadog.NewClientPlugin(ctx)
if err != nil {
    log.Fatal(err)
}
client, err := redis.NewClient(ctx, datadogPlugin.Register)

// For cluster client
datadogClusterPlugin, err := datadog.NewClusterPlugin(ctx)
if err != nil {
    log.Fatal(err)
}
cluster, err := redis.NewClusterClient(ctx, datadogClusterPlugin.Register)
```

### OpenTelemetry Plugin

```go
// For standalone client
otelPlugin, err := otel.NewClientPlugin(ctx)
if err != nil {
    log.Fatal(err)
}
client, err := redis.NewClient(ctx, otelPlugin.Register)

// For cluster client
otelClusterPlugin, err := otel.NewClusterPlugin(ctx)
if err != nil {
    log.Fatal(err)
}
cluster, err := redis.NewClusterClient(ctx, otelClusterPlugin.Register)
```

## Example: Complete Redis Client with Health Check

```go
package main

import (
    "context"
    "encoding/json"
    
    "github.com/xgodev/boost"
    h "github.com/xgodev/boost/extra/health"
    "github.com/xgodev/boost/factory/contrib/redis/go-redis/v9"
    "github.com/xgodev/boost/factory/contrib/redis/go-redis/v9/plugins/local/extra/health"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Initialize Boost
    boost.Start()
    
    // Create health check integrator
    healthIntegrator, err := health.NewClientHealth()
    if err != nil {
        log.Fatalf(err.Error())
    }
    
    // Create Redis client with health check
    client, err := redis.NewClient(context.Background(), healthIntegrator.Register)
    if err != nil {
        log.Error(err)
    }
    
    // Check health status
    all := h.CheckAll(context.Background())
    
    // Output health status
    j, _ := json.Marshal(all)
    log.Info(string(j))
}
```

## Integration with Boost

The Redis integration works seamlessly with other Boost components:

- **Logging**: Automatically integrates with Boost's logging system
- **Configuration**: Uses Boost's configuration management
- **Health Checks**: Built-in health check implementations
- **Metrics**: Can be integrated with Prometheus for metrics collection
- **Tracing**: Supports Datadog and OpenTelemetry for distributed tracing

## Best Practices

1. **Use Configuration Files**: Store Redis connection settings in external configuration files
2. **Enable Health Checks**: Always integrate health checks for monitoring Redis availability
3. **Add Observability**: Use Prometheus metrics and tracing plugins in production environments
4. **Handle Connection Errors**: Implement proper error handling for Redis connection issues
5. **Set Appropriate Timeouts**: Configure connection, read, and write timeouts based on your use case
6. **Pool Management**: Adjust pool size and idle connections based on your application's needs
7. **Use Cluster for High Availability**: Consider Redis Cluster for production environments requiring high availability

## Contributing

Contributions to improve the Redis integration are welcome. Please follow the Boost project's contribution guidelines.

## License

This package is part of the Boost project and is subject to its license terms.
