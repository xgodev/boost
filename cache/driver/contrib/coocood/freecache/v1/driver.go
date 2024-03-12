package freecache

import (
	"context"

	"github.com/coocood/freecache"
)

type Driver struct {
	cache   *freecache.Cache
	options *Options
}

func (c *Driver) Del(ctx context.Context, key string) error {
	c.cache.Del([]byte(key))
	return nil
}

func (c *Driver) Get(ctx context.Context, key string) (data []byte, err error) {
	data, err = c.cache.Get([]byte(key))
	if err != nil && err.Error() == "Entry not found" {
		return []byte{}, nil
	}
	return data, err
}

func (c *Driver) Set(ctx context.Context, key string, data []byte) (err error) {

	seconds := c.options.TTL.Seconds()

	if err = c.cache.Set([]byte(key), data, int(seconds)); err != nil {
		return err
	}

	return nil
}

func New(cache *freecache.Cache, options *Options) *Driver {
	return &Driver{cache: cache, options: options}
}
