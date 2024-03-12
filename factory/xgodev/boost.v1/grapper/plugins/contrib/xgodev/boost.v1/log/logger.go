package log

import (
	"context"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/grapper"
	"github.com/xgodev/boost/grapper/middleware/contrib/americanas-go/log.v1"
)

func NewAnyError[R any](ctx context.Context, name string) grapper.AnyErrorMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return log.NewAnyErrorMiddleware[R](ctx)
}

func NewAny[R any](ctx context.Context, name string) grapper.AnyMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return log.NewAnyMiddleware[R](ctx)
}

func NewError(ctx context.Context, name string) grapper.ErrorMiddleware {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return log.NewErrorMiddleware(ctx)
}
