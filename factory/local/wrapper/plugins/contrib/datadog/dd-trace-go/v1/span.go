package datadog

import (
	"context"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/wrapper"
	"github.com/xgodev/boost/wrapper/middleware/contrib/datadog/dd-trace-go/v1"
)

func NewAnyError[R any](ctx context.Context, name string) wrapper.AnyErrorMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return datadog.datadog.NewAnyErrorMiddleware[R](ctx, name, "wrapper")
}

func NewAny[R any](ctx context.Context, name string) wrapper.AnyMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return datadog.NewAnyMiddleware[R](ctx, name, "wrapper")
}

func NewError(ctx context.Context, name string) wrapper.ErrorMiddleware {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	return datadog.NewErrorMiddleware(ctx, name, "wrapper")
}
