package prometheus

import (
	"context"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/middleware"
	"github.com/xgodev/boost/middleware/plugins/contrib/prometheus/client_golang/v1"
)

func NewAnyError[R any](ctx context.Context, name string) middleware.AnyErrorMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return prometheus.NewAnyErrorMiddleware[R](ctx)
}

func NewAny[R any](ctx context.Context, name string) middleware.AnyMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return prometheus.NewAnyMiddleware[R](ctx)
}

func NewError(ctx context.Context, name string) middleware.ErrorMiddleware {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return prometheus.NewErrorMiddleware(ctx)
}
