package cache

import (
	"context"

	"github.com/xgodev/boost/cache"
	"github.com/xgodev/boost/wrapper"
)

type anyErrorMiddleware[R any] struct {
	manager *cache.Manager[R]
	opts    []cache.OptionSet
}

func (m *anyErrorMiddleware[R]) Exec(c *wrapper.AnyErrorContext[R], exec wrapper.AnyErrorExecFunc[R], returnFunc wrapper.AnyErrorReturnFunc[R]) (R, error) {
	return m.manager.GetOrSet(c.GetContext(), c.GetID(), func(ctx context.Context) (R, error) {
		return c.Next(exec, returnFunc)
	}, m.opts...)
}

func NewAnyErrorMiddleware[R any](ctx context.Context, manager *cache.Manager[R], opts ...cache.OptionSet) wrapper.AnyErrorMiddleware[R] {
	return &anyErrorMiddleware[R]{manager: manager, opts: opts}
}

type anyMiddleware[R any] struct {
	manager *cache.Manager[R]
	opts    []cache.OptionSet
}

func (m *anyMiddleware[R]) Exec(c *wrapper.AnyContext[R], exec wrapper.AnyExecFunc[R], returnFunc wrapper.AnyReturnFunc[R]) R {
	d, _ := m.manager.GetOrSet(c.GetContext(), c.GetID(), func(ctx context.Context) (R, error) {
		return c.Next(exec, returnFunc), nil
	}, m.opts...)
	return d
}

func NewAnyMiddleware[R any](ctx context.Context, manager *cache.Manager[R], opts ...cache.OptionSet) wrapper.AnyMiddleware[R] {
	return &anyMiddleware[R]{manager: manager, opts: opts}
}
