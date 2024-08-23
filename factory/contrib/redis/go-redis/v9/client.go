package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/xgodev/boost/wrapper/log"
)

// Plugin defines a function to process plugin.
type Plugin func(context.Context, *redis.Client) error

// NewClient returns a new redis client.
func NewClient(ctx context.Context, plugins ...Plugin) (*redis.Client, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithConfigPath returns a new redis client with options from config path.
func NewClientWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*redis.Client, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, opts, plugins...)
}

// NewClientWithOptions returns a new redis client with options.
func NewClientWithOptions(ctx context.Context, o *Options, plugins ...Plugin) (client *redis.Client, err error) {

	logger := log.FromContext(ctx)

	if redisSentinel(o) {
		client = failOverClient(o)
	} else {
		client = standaloneClient(o)
	}

	ping := client.Conn().Ping(ctx)
	if ping.Err() != nil {
		return nil, ping.Err()
	}

	for _, plugin := range plugins {
		if err := plugin(ctx, client); err != nil {
			return nil, err
		}
	}

	logger.Infof("Connected to Redis server: %s %s", client.Options().Addr, ping.String())

	return client, err
}

func failOverClient(o *Options) *redis.Client {
	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       o.Sentinel.MasterName,
		SentinelAddrs:    o.Sentinel.Addrs,
		SentinelPassword: o.Sentinel.Password,
		Password:         o.Password,
		MaxRetries:       o.MaxRetries,
		MinRetryBackoff:  o.MinRetryBackoff,
		MaxRetryBackoff:  o.MaxRetryBackoff,
		DialTimeout:      o.DialTimeout,
		DB:               o.Client.DB,
		ReadTimeout:      o.ReadTimeout,
		WriteTimeout:     o.WriteTimeout,
		PoolSize:         o.PoolSize,
		MinIdleConns:     o.MinIdleConns,
		PoolTimeout:      o.PoolTimeout,
	})
}

func standaloneClient(o *Options) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:            o.Client.Addr,
		Network:         o.Client.Network,
		Password:        o.Password,
		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,
		DialTimeout:     o.DialTimeout,
		DB:              o.Client.DB,
		ReadTimeout:     o.ReadTimeout,
		WriteTimeout:    o.WriteTimeout,
		PoolSize:        o.PoolSize,
		MinIdleConns:    o.MinIdleConns,
		PoolTimeout:     o.PoolTimeout,
	})
}

func redisSentinel(o *Options) bool {
	return o.Sentinel.MasterName != "" || o.Sentinel.Addrs != nil
}
