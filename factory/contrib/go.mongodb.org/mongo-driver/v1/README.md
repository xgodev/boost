# MongoDB Driver Factory - Advanced Configuration

This document provides detailed information about the advanced configuration options available for the MongoDB driver factory in the Boost framework.

## Overview

The MongoDB driver factory provides a comprehensive set of configuration options that allow you to fine-tune your MongoDB connections for various production scenarios. These options cover connection pooling, timeouts, authentication, TLS/SSL, read/write concerns, and more.

## Configuration Options

All configuration options are accessible through the Boost configuration system. You can set these options programmatically or through configuration files.

### Basic Connection

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.uri` | `mongodb://localhost:27017/temp` | MongoDB connection URI |
| `boost.factory.mongo.auth.username` | `""` | MongoDB username |
| `boost.factory.mongo.auth.password` | `""` | MongoDB password |

### Connection Pool

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.pool.max_size` | `100` | Maximum number of connections in the pool |
| `boost.factory.mongo.pool.min_size` | `0` | Minimum number of connections in the pool |
| `boost.factory.mongo.pool.max_idle_time_ms` | `0` | Maximum idle time for a pooled connection in milliseconds (0 = no limit) |

### Timeouts

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.timeout.connect_ms` | `30000` | Timeout for initial connection in milliseconds |
| `boost.factory.mongo.timeout.socket_ms` | `0` | Timeout for socket operations in milliseconds (0 = no timeout) |
| `boost.factory.mongo.timeout.server_selection_ms` | `30000` | Timeout for server selection in milliseconds |
| `boost.factory.mongo.timeout.heartbeat_ms` | `10000` | Interval between server monitoring checks in milliseconds |
| `boost.factory.mongo.timeout.local_threshold_ms` | `15` | Maximum latency difference between fastest and acceptable servers |
| `boost.factory.mongo.timeout.max_connecting` | `2` | Maximum number of concurrent connection attempts |
| `boost.factory.mongo.timeout.disable_ocsp_endpoint_check` | `false` | Disable OCSP endpoint checking |

### TLS/SSL

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.tls.enabled` | `false` | Enable TLS/SSL for MongoDB connection |
| `boost.factory.mongo.tls.insecure` | `false` | Allow insecure TLS/SSL connections (skip verification) |
| `boost.factory.mongo.tls.certificate_key_file` | `""` | Path to the client certificate and private key file |
| `boost.factory.mongo.tls.certificate_password` | `""` | Password for the client certificate private key |
| `boost.factory.mongo.tls.ca_file` | `""` | Path to the CA certificate file |

### Read Concern and Preference

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.read_concern.level` | `""` | Read concern level (available, local, majority, linearizable, snapshot) |
| `boost.factory.mongo.read_preference.mode` | `"primary"` | Read preference mode (primary, primaryPreferred, secondary, secondaryPreferred, nearest) |
| `boost.factory.mongo.read_preference.max_staleness_ms` | `90000` | Maximum staleness for secondary reads in milliseconds |

### Write Concern

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.write_concern.level` | `""` | Write concern level (majority, etc.) |
| `boost.factory.mongo.write_concern.w` | `1` | Write concern w value (number of nodes that must acknowledge writes) |
| `boost.factory.mongo.write_concern.j` | `false` | Write concern j value (whether writes should be journaled) |
| `boost.factory.mongo.write_concern.wtimeout_ms` | `0` | Write concern timeout in milliseconds (0 = no timeout) |

### Compression

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.compression.compressors` | `[]` | Compression algorithms to use (snappy, zlib, zstd) |
| `boost.factory.mongo.compression.zlib_level` | `6` | Compression level for zlib (0-9) |
| `boost.factory.mongo.compression.zstd_level` | `6` | Compression level for zstd (1-20) |

### Retry Configuration

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.retry.reads` | `true` | Enable retryable reads |
| `boost.factory.mongo.retry.writes` | `true` | Enable retryable writes |

### Replica Set

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.replica_set.name` | `""` | Replica set name |
| `boost.factory.mongo.replica_set.direct` | `false` | Use direct connection (bypass mongos) |

### Server API

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.server_api.version` | `""` | Server API version (1) |
| `boost.factory.mongo.server_api.strict` | `false` | Enable strict server API mode |
| `boost.factory.mongo.server_api.deprecation_errors` | `false` | Treat deprecated server API errors as errors |

### Miscellaneous

| Option | Default | Description |
|--------|---------|-------------|
| `boost.factory.mongo.app_name` | `""` | Application name for MongoDB logs and profiling |
| `boost.factory.mongo.load_balanced` | `false` | Enable load balanced mode |

## Configuration Priority

The MongoDB driver factory follows this configuration priority order:

1. **URI Parameters**: Settings specified in the connection URI have the highest priority
2. **Explicit Options**: Settings specified through configuration options have lower priority

This means that any parameter specified in the URI will override the same parameter set through configuration options.

## Usage Examples

### Basic Connection with Advanced Options

```go
package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {
	// Initialize Boost framework
	boost.Start()
	
	// Configure MongoDB options
	config.Set("boost.factory.mongo.pool.max_size", 100)
	config.Set("boost.factory.mongo.pool.min_size", 10)
	config.Set("boost.factory.mongo.timeout.connect_ms", 5000)
	config.Set("boost.factory.mongo.read_preference.mode", "secondaryPreferred")
	config.Set("boost.factory.mongo.write_concern.j", true)
	
	// URI parameters will override any conflicting options above
	config.Set("boost.factory.mongo.uri", "mongodb://localhost:27017/mydb?maxPoolSize=50&readPreference=primary")
	
	// Create MongoDB connection
	conn, err := mongo.NewConn(context.Background())
	if err != nil {
		log.Errorf("Failed to connect: %v", err)
		return
	}
	
	// Use the connection...
	
	// Close the connection when done
	conn.Client.Disconnect(context.Background())
}
```

### Secure Connection with TLS

```go
package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {
	// Initialize Boost framework
	boost.Start()
	
	// Configure MongoDB with TLS
	config.Set("boost.factory.mongo.uri", "mongodb://localhost:27017/mydb")
	config.Set("boost.factory.mongo.tls.enabled", true)
	config.Set("boost.factory.mongo.tls.ca_file", "/path/to/ca.pem")
	config.Set("boost.factory.mongo.tls.certificate_key_file", "/path/to/client.pem")
	
	// Create MongoDB connection
	conn, err := mongo.NewConn(context.Background())
	if err != nil {
		log.Errorf("Failed to connect: %v", err)
		return
	}
	
	// Use the connection...
	
	// Close the connection when done
	conn.Client.Disconnect(context.Background())
}
```

### High Availability Configuration

```go
package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {
	// Initialize Boost framework
	boost.Start()
	
	// Configure MongoDB for high availability
	config.Set("boost.factory.mongo.uri", "mongodb://server1:27017,server2:27017,server3:27017/mydb")
	config.Set("boost.factory.mongo.replica_set.name", "myReplicaSet")
	config.Set("boost.factory.mongo.read_preference.mode", "nearest")
	config.Set("boost.factory.mongo.write_concern.level", "majority")
	config.Set("boost.factory.mongo.write_concern.j", true)
	config.Set("boost.factory.mongo.write_concern.wtimeout_ms", 5000)
	config.Set("boost.factory.mongo.retry.reads", true)
	config.Set("boost.factory.mongo.retry.writes", true)
	
	// Create MongoDB connection
	conn, err := mongo.NewConn(context.Background())
	if err != nil {
		log.Errorf("Failed to connect: %v", err)
		return
	}
	
	// Use the connection...
	
	// Close the connection when done
	conn.Client.Disconnect(context.Background())
}
```

## Best Practices

1. **Connection Pooling**: Configure pool sizes based on your application's concurrency requirements. For high-throughput applications, increase the `max_pool_size` and set an appropriate `min_pool_size` to maintain warm connections.

2. **Timeouts**: Set appropriate timeouts for your environment. In production, you may want shorter timeouts to fail fast and allow for retries.

3. **Read/Write Concerns**: Choose read and write concerns based on your consistency and availability requirements:
   - For strong consistency, use `majority` read concern and write concern
   - For higher performance with eventual consistency, use `local` read concern and lower write concern values

4. **Retries**: Enable retry reads and writes for better resilience in distributed environments.

5. **Monitoring**: Set a meaningful `app_name` to identify your application in MongoDB logs and monitoring tools.

6. **Security**: Always use TLS in production environments and consider using certificate authentication for sensitive deployments.

7. **Compression**: Enable compression for high-volume applications, especially over slower network connections.

## Troubleshooting

If you encounter connection issues:

1. Check that the MongoDB server is running and accessible
2. Verify that the connection URI is correct
3. Ensure that authentication credentials are valid
4. Check TLS configuration if enabled
5. Verify that timeout values are appropriate for your network environment
6. Check server logs for any connection errors

For performance issues:

1. Review connection pool settings
2. Check read/write concern configurations
3. Monitor connection pool utilization
4. Consider enabling compression for large data transfers
5. Review server selection timeout and local threshold settings

## Further Reading

- [MongoDB Driver Documentation](https://docs.mongodb.com/drivers/go/)
- [MongoDB Connection String URI Format](https://docs.mongodb.com/manual/reference/connection-string/)
- [MongoDB Read Concerns](https://docs.mongodb.com/manual/reference/read-concern/)
- [MongoDB Write Concerns](https://docs.mongodb.com/manual/reference/write-concern/)
