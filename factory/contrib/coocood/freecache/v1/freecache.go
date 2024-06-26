package freecache

import (
	"context"

	"github.com/coocood/freecache"
	"github.com/xgodev/boost/wrapper/log"
)

// NewCacheWithConfigPath returns a cache with options from config path .
func NewCacheWithConfigPath(ctx context.Context, path string, opts ...Option) (*freecache.Cache, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewCacheWithOptions(ctx, options, opts...)
}

// NewCacheWithOptions returns a cache with options.
func NewCacheWithOptions(ctx context.Context, o *Options, opts ...Option) (cache *freecache.Cache, err error) {

	logger := log.FromContext(ctx)

	for _, opt := range opts {
		opt(o)
	}

	cache = freecache.NewCache(o.CacheSize)

	logger.Infof("Created cache with size %v", o.CacheSize)

	return cache, err
}

// NewCache returns a cache.
func NewCache(ctx context.Context, opts ...Option) (*freecache.Cache, error) {

	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewCacheWithOptions(ctx, o, opts...)
}
