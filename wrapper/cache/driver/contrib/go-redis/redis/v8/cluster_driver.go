package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Cluster struct {
	cache   *redis.ClusterClient
	options *Options
}

func (c *Cluster) Del(ctx context.Context, key string) error {
	return c.cache.Del(ctx, key).Err()
}

func (c *Cluster) Get(ctx context.Context, key string) (data []byte, err error) {
	return c.cache.Get(ctx, key).Bytes()
}

func (c *Cluster) Set(ctx context.Context, key string, data []byte) (err error) {

	c.cache.Del(ctx, key)

	if err = c.cache.Set(ctx, key, data, c.options.TTL).Err(); err != nil {
		return err
	}

	return nil
}

func NewCluster(cache *redis.ClusterClient, options *Options) *Cluster {
	return &Cluster{cache: cache, options: options}
}
