# Google Cloud Pub/Sub Factory

## Overview

The Google Cloud Pub/Sub Factory provides integration with the [Google Cloud Pub/Sub](https://cloud.google.com/pubsub) service within the Boost framework. This factory enables seamless messaging and event streaming capabilities for applications built on Boost, allowing for reliable, scalable, and secure communication between services.

## Features

- **Simplified Client Creation**: Easy creation of Pub/Sub clients with sensible defaults
- **Flexible Configuration**: Comprehensive configuration options through Boost's configuration system
- **Plugin Support**: Extensible plugin architecture for customizing client behavior
- **Integration with GCP Auth**: Seamless authentication with Google Cloud services
- **gRPC Optimization**: Fine-tuned gRPC settings for optimal performance
- **Logging Integration**: Automatic integration with Boost's logging system

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Initialize Boost
    boost.Start()
    
    // Create a context with logger
    ctx := log.WithLogger(context.Background(), log.GetLogger())
    
    // Create a Pub/Sub client with default settings
    client, err := pubsub.NewClient(ctx)
    if err != nil {
        log.Errorf("Failed to create Pub/Sub client: %v", err)
        return
    }
    defer client.Close()
    
    // Use the client
    topic := client.Topic("my-topic")
    result := topic.Publish(ctx, &pubsub.Message{
        Data: []byte("Hello, Pub/Sub!"),
    })
    
    // Get the server-generated ID for the published message
    id, err := result.Get(ctx)
    if err != nil {
        log.Errorf("Failed to publish message: %v", err)
        return
    }
    
    log.Infof("Published message with ID: %s", id)
}
```

### Custom Configuration Path

```go
// Create a client with custom configuration path
client, err := pubsub.NewClientWithConfigPath(ctx, "myapp.pubsub")
if err != nil {
    log.Errorf("Failed to create Pub/Sub client: %v", err)
    return
}
defer client.Close()
```

### Custom Options

```go
// Create options and modify them
options, err := pubsub.NewOptions()
if err != nil {
    log.Errorf("Failed to create options: %v", err)
    return
}

// Customize options
options.APIOptions.ProjectID = "my-gcp-project"
options.GRPCOptions.MaxRecvMsgSize = 10 * 1024 * 1024 // 10 MB

// Create client with custom options
client, err := pubsub.NewClientWithOptions(ctx, options)
if err != nil {
    log.Errorf("Failed to create Pub/Sub client: %v", err)
    return
}
defer client.Close()
```

### Using Plugins

```go
// Import the gRPC client plugin package
import clientgrpc "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"

// Define a custom plugin
myPlugin := func(ctx context.Context, opts ...interface{}) (interface{}, error) {
    // Custom logic for the plugin
    return nil, nil
}

// Create client with plugin
client, err := pubsub.NewClient(ctx, myPlugin)
if err != nil {
    log.Errorf("Failed to create Pub/Sub client: %v", err)
    return
}
defer client.Close()
```

## Configuration Parameters

The Pub/Sub factory uses a combination of API options and gRPC options for configuration.

### API Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `apiOptions.projectID` | Google Cloud project ID | From environment or application default credentials |
| `apiOptions.credentialsFile` | Path to service account credentials file | "" (uses application default credentials) |
| `apiOptions.audience` | Target audience for authentication | "" |
| `apiOptions.quotaProject` | Project for quota and billing | "" |
| `apiOptions.requestReason` | Reason for the request (for auditing) | "" |
| `apiOptions.scopes` | OAuth scopes for authentication | Default Pub/Sub scopes |
| `apiOptions.endpoint` | API endpoint override | "" (uses default endpoint) |

### gRPC Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `grpcOptions.maxRecvMsgSize` | Maximum message size client can receive | 4 MB |
| `grpcOptions.maxSendMsgSize` | Maximum message size client can send | 4 MB |
| `grpcOptions.keepaliveTime` | After this duration of no activity, ping server to check if connection is alive | 10s |
| `grpcOptions.keepaliveTimeout` | Time to wait for ping ack before considering connection dead | 20s |
| `grpcOptions.keepalivePermitWithoutStream` | Allow pings even without active streams | true |
| `grpcOptions.initialWindowSize` | Initial window size for stream flow control | 0 (system default) |
| `grpcOptions.initialConnWindowSize` | Initial window size for connection flow control | 0 (system default) |
| `grpcOptions.maxRetryRPCBufferSize` | Maximum buffer size for retry RPCs | 0 (system default) |
| `grpcOptions.maxConnectionIdle` | Maximum time a connection can be idle | 0 (infinite) |
| `grpcOptions.maxConnectionAge` | Maximum age of a connection | 0 (infinite) |
| `grpcOptions.maxConnectionAgeGrace` | Grace period after max connection age | 0 (infinite) |

## Integration with Other Boost Components

The Google Cloud Pub/Sub Factory integrates with:

- **Config Wrapper**: For loading and managing configuration
- **Log Wrapper**: For logging Pub/Sub operations and errors
- **Google Cloud API Factory**: For shared API configuration
- **gRPC Factory**: For optimizing gRPC connections

## Example: Publishing and Subscribing

```go
package main

import (
    "context"
    "fmt"
    "sync"
    
    "cloud.google.com/go/pubsub"
    "github.com/xgodev/boost"
    pubsubfactory "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    boost.Start()
    
    ctx := log.WithLogger(context.Background(), log.GetLogger())
    
    // Create a Pub/Sub client
    client, err := pubsubfactory.NewClient(ctx)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()
    
    // Create a topic if it doesn't exist
    topicID := "my-topic"
    topic := client.Topic(topicID)
    exists, err := topic.Exists(ctx)
    if err != nil {
        log.Fatalf("Failed to check if topic exists: %v", err)
    }
    
    if !exists {
        topic, err = client.CreateTopic(ctx, topicID)
        if err != nil {
            log.Fatalf("Failed to create topic: %v", err)
        }
        log.Infof("Topic created: %s", topicID)
    }
    
    // Create a subscription if it doesn't exist
    subID := "my-subscription"
    sub := client.Subscription(subID)
    exists, err = sub.Exists(ctx)
    if err != nil {
        log.Fatalf("Failed to check if subscription exists: %v", err)
    }
    
    if !exists {
        sub, err = client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
            Topic:               topic,
            AckDeadline:         20, // seconds
            RetainAckedMessages: false,
        })
        if err != nil {
            log.Fatalf("Failed to create subscription: %v", err)
        }
        log.Infof("Subscription created: %s", subID)
    }
    
    // Publish a message
    result := topic.Publish(ctx, &pubsub.Message{
        Data: []byte("Hello, Pub/Sub!"),
        Attributes: map[string]string{
            "origin": "boost-example",
        },
    })
    
    // Get the server-generated ID for the published message
    id, err := result.Get(ctx)
    if err != nil {
        log.Fatalf("Failed to publish message: %v", err)
    }
    log.Infof("Published message with ID: %s", id)
    
    // Receive messages
    var mu sync.Mutex
    received := 0
    
    err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
        mu.Lock()
        defer mu.Unlock()
        log.Infof("Received message: %s", string(msg.Data))
        log.Infof("Attributes: %v", msg.Attributes)
        received++
        msg.Ack()
        
        if received >= 1 {
            // Cancel context to stop receiving after getting our message
            cancel, ok := ctx.Value("cancel").(context.CancelFunc)
            if ok {
                cancel()
            }
        }
    })
    
    if err != nil && err != context.Canceled {
        log.Fatalf("Failed to receive messages: %v", err)
    }
}
```

## Best Practices

1. **Project ID**: Always specify your Google Cloud project ID in production environments
2. **Authentication**: Use service account credentials with the minimum required permissions
3. **Error Handling**: Implement robust error handling for Pub/Sub operations
4. **Message Size**: Be mindful of message size limits (default 10MB)
5. **Subscription Management**: Set appropriate ack deadlines based on your processing time
6. **Retry Logic**: Implement retry logic for transient failures
7. **Monitoring**: Set up monitoring for Pub/Sub topics and subscriptions

## References

- [Google Cloud Pub/Sub Documentation](https://cloud.google.com/pubsub/docs)
- [Go Client for Google Cloud Pub/Sub](https://pkg.go.dev/cloud.google.com/go/pubsub)
- [Boost Framework Documentation](https://github.com/xgodev/boost)
