package cache

import (
	"context"

	"github.com/xgodev/boost/middleware"
	"github.com/xgodev/boost/wrapper/cache"
)

type anyErrorMiddleware[R any] struct {
	manager *cache.Manager[R]
	opts    []cache.OptionSet
}

func (m *anyErrorMiddleware[R]) Exec(c *middleware.AnyErrorContext[R], exec middleware.AnyErrorExecFunc[R], returnFunc middleware.AnyErrorReturnFunc[R]) (R, error) {
	return m.manager.GetOrSet(c.GetContext(), c.GetID(), func(ctx context.Context) (R, error) {
		return c.Next(exec, returnFunc)
	}, m.opts...)
}

func NewAnyErrorMiddleware[R any](ctx context.Context, manager *cache.Manager[R], opts ...cache.OptionSet) middleware.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{manager: manager, opts: opts}
}

type anyMiddleware[R any] struct {
	manager *cache.Manager[R]
	opts    []cache.OptionSet
}

func (m *anyMiddleware[R]) Exec(c *middleware.AnyContext[R], exec middleware.AnyExecFunc[R], returnFunc middleware.AnyReturnFunc[R]) R {
	d, _ := m.manager.GetOrSet(c.GetContext(), c.GetID(), func(ctx context.Context) (R, error) {
		return c.Next(exec, returnFunc), nil
	}, m.opts...)
	return d
}

func NewAnyMiddleware[R any](ctx context.Context, manager *cache.Manager[R], opts ...cache.OptionSet) middleware.AnyMiddleware[R] {
	return &anyMiddleware[R]{manager: manager, opts: opts}
}
