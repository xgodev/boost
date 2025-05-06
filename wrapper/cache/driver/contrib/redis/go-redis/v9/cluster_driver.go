package redis

import (
	"context"
	"github.com/xgodev/boost/model/errors"

	"github.com/redis/go-redis/v9"
)

type Cluster struct {
	cache   *redis.ClusterClient
	options *Options
}

// Del removes the key from the cluster.
func (c *Cluster) Del(ctx context.Context, key string) error {
	return c.cache.Del(ctx, key).Err()
}

// Get performs a single GET and treats redis.Nil as cache miss.
func (c *Cluster) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.cache.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// Set writes the key with TTL, without a prior DEL.
func (c *Cluster) Set(ctx context.Context, key string, data []byte) error {
	return c.cache.Set(ctx, key, data, c.options.TTL).Err()
}

// NewCluster constructs a new Cluster driver, ensuring options is non-nil.
func NewCluster(cache *redis.ClusterClient, options *Options) *Cluster {
	if options == nil {
		options = &Options{TTL: 0}
	}
	return &Cluster{cache: cache, options: options}
}
