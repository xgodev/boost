package datadog

import (
	"context"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/extra/middleware/plugins/contrib/datadog/dd-trace-go/v1"
	"github.com/xgodev/boost/wrapper/config"
)

func NewAnyError[R any](ctx context.Context, name string) middleware.AnyErrorMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return datadog.NewAnyErrorMiddleware[R](ctx, name, "wrapper")
}

func NewAny[R any](ctx context.Context, name string) middleware.AnyMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return datadog.NewAnyMiddleware[R](ctx, name, "wrapper")
}

func NewError(ctx context.Context, name string) middleware.ErrorMiddleware {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return datadog.NewErrorMiddleware(ctx, name, "wrapper")
}
