package log

import (
	"context"
	"github.com/xgodev/boost/middleware"
	"github.com/xgodev/boost/middleware/plugins/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/config"
)

func NewAnyError[R any](ctx context.Context, name string) middleware.AnyErrorMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return log.NewAnyErrorMiddleware[R](ctx)
}

func NewAny[R any](ctx context.Context, name string) middleware.AnyMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return log.NewAnyMiddleware[R](ctx)
}

func NewError(ctx context.Context, name string) middleware.ErrorMiddleware {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return log.NewErrorMiddleware(ctx)
}
