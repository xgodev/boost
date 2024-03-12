package redis

import (
	"context"

	"github.com/go-redis/redis/v7"
)

type Client struct {
	cache   *redis.Client
	options *Options
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.cache.WithContext(ctx).Del(key).Err()
}

func (c *Client) Get(ctx context.Context, key string) (data []byte, err error) {
	return c.cache.WithContext(ctx).Get(key).Bytes()
}

func (c *Client) Set(ctx context.Context, key string, data []byte) (err error) {
	c.cache.WithContext(ctx).Del(key)

	if err = c.cache.WithContext(ctx).Set(key, data, c.options.TTL).Err(); err != nil {
		return err
	}

	return nil
}

func NewClient(cache *redis.Client, options *Options) *Client {
	return &Client{cache: cache, options: options}
}
