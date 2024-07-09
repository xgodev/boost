package cache

import (
	"context"
	"github.com/xgodev/boost/middleware"
	gcache "github.com/xgodev/boost/middleware/plugins/local/wrapper/cache"
	"github.com/xgodev/boost/wrapper/cache"
	"github.com/xgodev/boost/wrapper/config"
)

type Cache[T any] struct {
	manager *cache.Manager[T]
}

func New[T any](ctx context.Context, manager *cache.Manager[T]) *Cache[T] {
	return &Cache[T]{manager: manager}
}

func (c *Cache[T]) NewAnyError(ctx context.Context, name string) middleware.AnyErrorMiddleware[T] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return gcache.NewAnyErrorMiddleware[T](ctx, c.manager)
}

func (c *Cache[T]) NewAny(ctx context.Context, name string) middleware.AnyMiddleware[T] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return gcache.NewAnyMiddleware[T](ctx, c.manager)
}
