package hystrix

import (
	"context"
	"github.com/xgodev/boost/config"
	h "github.com/xgodev/boost/factory/afex/hystrix-go.v0"
	"github.com/xgodev/boost/grapper"
	"github.com/xgodev/boost/grapper/middleware/contrib/afex/hystrix-go.v0"
	"github.com/xgodev/boost/log"
)

func NewAnyError[R any](ctx context.Context, name string) grapper.AnyErrorMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	if err := h.ConfigureCommand(name); err != nil {
		log.Error(err.Error())
	}
	return hystrix.NewAnyErrorMiddleware[R](ctx, name)
}

func NewAny[R any](ctx context.Context, name string) grapper.AnyMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	if err := h.ConfigureCommand(name); err != nil {
		log.Error(err.Error())
	}
	return hystrix.NewAnyMiddleware[R](ctx, name)
}

func NewError(ctx context.Context, name string) grapper.ErrorMiddleware {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	if err := h.ConfigureCommand(name); err != nil {
		log.Error(err.Error())
	}
	return hystrix.NewErrorMiddleware(ctx, name)
}
