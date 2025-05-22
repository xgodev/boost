# AWS SDK Factory

## Overview

The AWS SDK Factory provides integration with the [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2), enabling seamless interaction with Amazon Web Services within the Boost framework. This factory simplifies AWS service configuration and usage through standardized interfaces and configuration management.

## Features

- **Unified Configuration**: Centralized AWS configuration management through Boost's configuration system
- **Credential Management**: Multiple credential sources including environment variables and configuration files
- **Custom Endpoints**: Support for custom endpoints to facilitate local development and testing
- **HTTP Client Customization**: Fine-grained control over HTTP client behavior
- **Retry Policies**: Configurable retry mechanisms with rate limiting options
- **Plugin System**: Extensible plugin architecture for customizing AWS client behavior

## Usage

### Basic Configuration

```go
package main

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/xgodev/boost"
    awsfactory "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Initialize Boost
    boost.Start()
    
    // Create a context with logger
    ctx := log.WithLogger(context.Background(), log.GetLogger())
    
    // Create AWS configuration
    cfg, err := awsfactory.NewConfig(ctx)
    if err != nil {
        log.Errorf("Failed to create AWS config: %v", err)
        return
    }
    
    // Create an AWS service client
    s3Client := s3.NewFromConfig(cfg)
    
    // Use the client
    // ...
}
```

### Custom Configuration Path

```go
// Create AWS configuration with custom configuration path
cfg, err := awsfactory.NewConfigWithConfigPath(ctx, "myapp.aws")
if err != nil {
    log.Errorf("Failed to create AWS config: %v", err)
    return
}
```

### Custom Options

```go
// Create options and modify them
options, err := awsfactory.NewOptions()
if err != nil {
    log.Errorf("Failed to create options: %v", err)
    return
}

// Customize options
options.DefaultRegion = "us-west-2"
options.MaxAttempts = 3
options.HasRateLimit = false

// Create AWS configuration with custom options
cfg, err := awsfactory.NewConfigWithOptions(ctx, options)
if err != nil {
    log.Errorf("Failed to create AWS config: %v", err)
    return
}
```

### Using Plugins

```go
// Define a custom plugin
myPlugin := func(ctx context.Context, cfg *aws.Config) error {
    // Customize AWS config
    cfg.Region = "eu-central-1"
    return nil
}

// Create AWS configuration with plugin
cfg, err := awsfactory.NewConfig(ctx, myPlugin)
if err != nil {
    log.Errorf("Failed to create AWS config: %v", err)
    return
}
```

### Custom Endpoints

Custom endpoints can be configured through Boost's configuration system:

```yaml
boost:
  factory:
    aws:
      customEndpoint:
        s3:
          url: "http://localhost:4566"
          signingRegion: "us-east-1"
```

## Configuration Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `accessKeyId` | AWS access key ID | Environment variable `AWS_ACCESS_KEY_ID` |
| `secretAccessKey` | AWS secret access key | Environment variable `AWS_SECRET_ACCESS_KEY` |
| `defaultRegion` | Default AWS region | Environment variable `AWS_DEFAULT_REGION` |
| `defaultAccountNumber` | Default AWS account number | Environment variable `AWS_DEFAULT_ACCOUNT_NUMBER` |
| `sessionToken` | AWS session token | Environment variable `AWS_SESSION_TOKEN` |
| `customEndpoint` | Map of service-specific endpoint configurations | Empty map |
| `retryer.maxAttempts` | Maximum number of retry attempts | 5 |
| `retryer.hasRateLimit` | Whether to use rate limiting for retries | true |

### HTTP Client Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `httpClient.maxIdleConnPerHost` | Maximum idle connections per host | 10 |
| `httpClient.maxIdleConn` | Maximum idle connections | 100 |
| `httpClient.maxConnsPerHost` | Maximum connections per host | 256 |
| `httpClient.idleConnTimeout` | Idle connection timeout | 90s |
| `httpClient.disableKeepAlives` | Disable HTTP keep-alives | true |
| `httpClient.disableCompression` | Disable HTTP compression | true |
| `httpClient.forceHTTP2` | Force HTTP/2 | true |
| `httpClient.TLSHandshakeTimeout` | TLS handshake timeout | 10s |
| `httpClient.timeout` | Request timeout | 30s |
| `httpClient.dialTimeout` | Dial timeout | 5s |
| `httpClient.keepAlive` | Keep-alive duration | 15s |
| `httpClient.expectContinueTimeout` | Expect-continue timeout | 1s |
| `httpClient.dualStack` | Use dual-stack addressing | true |

## Integration with Other Boost Components

The AWS SDK Factory integrates with:

- **Config Wrapper**: For loading and managing configuration
- **Log Wrapper**: For logging AWS operations and errors
- **HTTP Client Factory**: For customizing HTTP client behavior

## Plugin System

The factory supports plugins for extending functionality. Plugins are functions that receive and can modify the AWS configuration:

```go
type Plugin func(context.Context, *aws.Config) error
```

Common use cases for plugins include:

- Adding custom middleware to AWS clients
- Configuring service-specific options
- Setting up logging or metrics collection
- Implementing custom authentication mechanisms

## Best Practices

1. **Credential Management**: Avoid hardcoding credentials; use environment variables or IAM roles
2. **Region Selection**: Set appropriate regions based on your application's needs and user proximity
3. **Retry Configuration**: Adjust retry settings based on your application's requirements and AWS service limits
4. **HTTP Client Tuning**: Optimize HTTP client settings for your specific workload patterns
5. **Custom Endpoints**: Use custom endpoints for local development and testing with tools like LocalStack

## Example: S3 Client

```go
package main

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/xgodev/boost"
    awsfactory "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    boost.Start()
    
    ctx := log.WithLogger(context.Background(), log.GetLogger())
    
    // Create AWS configuration
    cfg, err := awsfactory.NewConfig(ctx)
    if err != nil {
        log.Errorf("Failed to create AWS config: %v", err)
        return
    }
    
    // Create S3 client
    s3Client := s3.NewFromConfig(cfg)
    
    // List buckets
    result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
    if err != nil {
        log.Errorf("Failed to list buckets: %v", err)
        return
    }
    
    // Process results
    for _, bucket := range result.Buckets {
        log.Infof("Bucket: %s, Created: %s", *bucket.Name, bucket.CreationDate)
    }
}
```

## References

- [AWS SDK for Go v2 GitHub Repository](https://github.com/aws/aws-sdk-go-v2)
- [AWS SDK for Go v2 Documentation](https://aws.github.io/aws-sdk-go-v2/docs/)
- [Boost Framework Documentation](https://github.com/xgodev/boost)
