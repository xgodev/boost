package log

import (
	"context"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/wrapper"
	"github.com/xgodev/boost/wrapper/middleware/local/log"
)

func NewAnyError[R any](ctx context.Context, name string) wrapper.AnyErrorMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return log.NewAnyErrorMiddleware[R](ctx)
}

func NewAny[R any](ctx context.Context, name string) wrapper.AnyMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return log.NewAnyMiddleware[R](ctx)
}

func NewError(ctx context.Context, name string) wrapper.ErrorMiddleware {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return log.NewErrorMiddleware(ctx)
}
