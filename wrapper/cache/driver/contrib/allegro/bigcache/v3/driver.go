package bigcache

import (
	"context"

	"github.com/allegro/bigcache/v3"
)

type Driver struct {
	cache *bigcache.BigCache
}

func (c *Driver) Del(ctx context.Context, key string) (err error) {
	err = c.cache.Delete(key)
	if err != nil && err.Error() == "Entry not found" {
		return nil
	}
	return err
}

func (c *Driver) Get(ctx context.Context, key string) (data []byte, err error) {
	data, err = c.cache.Get(key)
	if err != nil && err.Error() == "Entry not found" {
		return []byte{}, nil
	}
	return data, err
}

func (c *Driver) Set(ctx context.Context, key string, data []byte) (err error) {

	if err = c.cache.Set(key, data); err != nil {
		return err
	}

	return nil
}

func New(cache *bigcache.BigCache) *Driver {
	return &Driver{cache: cache}
}
