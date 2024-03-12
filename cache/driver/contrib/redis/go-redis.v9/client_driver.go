package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	cache   *redis.Client
	options *Options
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.cache.Del(ctx, key).Err()
}

func (c *Client) Get(ctx context.Context, key string) (data []byte, err error) {
	return c.cache.Get(ctx, key).Bytes()
}

func (c *Client) Set(ctx context.Context, key string, data []byte) (err error) {

	c.cache.Del(ctx, key)

	if err = c.cache.Set(ctx, key, data, c.options.TTL).Err(); err != nil {
		return err
	}

	return nil
}

func NewClient(cache *redis.Client, options *Options) *Client {
	return &Client{cache: cache, options: options}
}
