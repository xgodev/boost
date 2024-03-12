package datadog

import (
	"context"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/grapper"
	"github.com/xgodev/boost/grapper/middleware/contrib/datadog/dd-trace-go.v1"
)

func NewAnyError[R any](ctx context.Context, name string) grapper.AnyErrorMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return datadog.NewAnyErrorMiddleware[R](ctx, name, "wrapper")
}

func NewAny[R any](ctx context.Context, name string) grapper.AnyMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return datadog.NewAnyMiddleware[R](ctx, name, "wrapper")
}

func NewError(ctx context.Context, name string) grapper.ErrorMiddleware {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return datadog.NewErrorMiddleware(ctx, name, "wrapper")
}
