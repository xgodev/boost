package redis

import (
	"context"
	"github.com/xgodev/boost/model/errors"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	cache   *redis.Client
	options *Options
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.cache.Del(ctx, key).Err()
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.cache.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}
func (c *Client) Set(ctx context.Context, key string, data []byte) (err error) {
	return c.cache.Set(ctx, key, data, c.options.TTL).Err()
}

func NewClient(cache *redis.Client, options *Options) *Client {
	return &Client{cache: cache, options: options}
}
