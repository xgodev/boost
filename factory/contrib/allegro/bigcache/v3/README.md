# BigCache Factory

## Overview

The BigCache Factory provides integration with the [BigCache](https://github.com/allegro/bigcache) library, a fast, concurrent, evicting in-memory cache designed for high-load applications. This factory enables efficient caching capabilities within the Boost framework, optimized for applications that require high performance with minimal GC overhead.

## Features

- **In-Memory Caching**: Fast access to cached data with minimal latency
- **Configurable Eviction**: Time-based eviction policies with customizable life windows
- **Sharded Architecture**: Concurrent access with minimal lock contention
- **Memory Management**: Hard limits on cache size to prevent memory issues
- **Metrics & Statistics**: Optional statistics collection for monitoring
- **Logging Integration**: Seamless integration with Boost's logging system

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "github.com/xgodev/boost"
    "github.com/xgodev/boost/factory/contrib/allegro/bigcache/v3"
    "github.com/xgodev/boost/wrapper/log"
)

func main() {
    // Initialize Boost
    boost.Start()
    
    // Create a context with logger
    ctx := log.WithLogger(context.Background(), log.GetLogger())
    
    // Create a new cache with default settings
    cache, err := bigcache.NewCache(ctx)
    if err != nil {
        log.Errorf("Failed to create cache: %v", err)
        return
    }
    
    // Use the cache
    cache.Set("key", []byte("value"))
    entry, err := cache.Get("key")
    if err != nil {
        log.Errorf("Failed to get entry: %v", err)
    }
    
    log.Infof("Retrieved value: %s", string(entry))
}
```

### Custom Configuration

You can customize BigCache parameters through Boost's configuration system:

```go
// Create a cache with custom configuration path
cache, err := bigcache.NewCacheWithConfigPath(ctx, "myapp.cache")
if err != nil {
    log.Errorf("Failed to create cache: %v", err)
    return
}
```

Or by directly providing options:

```go
// Create options and modify them
options, err := bigcache.NewOptions()
if err != nil {
    log.Errorf("Failed to create options: %v", err)
    return
}

// Customize options
options.Shards = 2048
options.LifeWindow = 10 * time.Minute
options.HardMaxCacheSize = 512 // MB

// Create cache with custom options
cache, err := bigcache.NewCacheWithOptions(ctx, options)
if err != nil {
    log.Errorf("Failed to create cache: %v", err)
    return
}
```

## Configuration Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `shards` | Number of cache shards (must be a power of two) | 1024 |
| `lifeWindow` | Time after which entry can be evicted | 5 minutes |
| `cleanWindow` | Interval between removing expired entries | 0 (disabled) |
| `maxEntriesInWindow` | Max number of entries in life window | 600,000 |
| `maxEntrySize` | Max size of entry in bytes | 1 MB |
| `verbose` | Verbose mode for memory allocation information | false |
| `hardMaxCacheSize` | Limit for cache size in MB (0 = unlimited) | 1 |
| `statsEnabled` | Enable statistics collection | false |

## Integration with Other Boost Components

The BigCache Factory integrates with:

- **Config Wrapper**: For loading and managing configuration
- **Log Wrapper**: For logging cache events and errors
- **Context**: For propagating logger and other values

## Best Practices

1. **Shard Count**: Set the number of shards based on your expected concurrency level. Higher concurrency requires more shards to reduce lock contention.

2. **Life Window**: Choose an appropriate life window based on your data freshness requirements. Shorter windows reduce memory usage but may cause more cache misses.

3. **Entry Size**: Set the max entry size based on your typical cached objects. Setting it too high wastes memory, while setting it too low causes errors when storing larger objects.

4. **Memory Management**: Use `hardMaxCacheSize` to prevent the cache from consuming too much memory, especially in memory-constrained environments.

5. **Clean Window**: Enable periodic cleanup (by setting a positive `cleanWindow` value) for long-running applications to ensure memory is reclaimed.

## Performance Considerations

- BigCache is optimized for high-load scenarios with minimal GC overhead
- The cache performs best with a large number of small entries
- For very large objects, consider using a different caching solution
- Monitor cache statistics in production to fine-tune parameters

## References

- [BigCache GitHub Repository](https://github.com/allegro/bigcache)
- [Boost Framework Documentation](https://github.com/xgodev/boost)
