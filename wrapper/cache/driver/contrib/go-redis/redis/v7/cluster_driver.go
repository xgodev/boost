package redis

import (
	"context"

	"github.com/go-redis/redis/v7"
)

type Cluster struct {
	cache   *redis.ClusterClient
	options *Options
}

func (c *Cluster) Del(ctx context.Context, key string) error {
	return c.cache.WithContext(ctx).Del(key).Err()
}

func (c *Cluster) Get(ctx context.Context, key string) (data []byte, err error) {
	return c.cache.WithContext(ctx).Get(key).Bytes()
}

func (c *Cluster) Set(ctx context.Context, key string, data []byte) (err error) {

	c.cache.WithContext(ctx).Del(key)

	if err = c.cache.WithContext(ctx).Set(key, data, c.options.TTL).Err(); err != nil {
		return err
	}

	return nil
}

func NewCluster(cache *redis.ClusterClient, options *Options) *Cluster {
	return &Cluster{cache: cache, options: options}
}
