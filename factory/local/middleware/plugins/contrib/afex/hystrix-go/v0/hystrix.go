package hystrix

import (
	"context"
	h "github.com/xgodev/boost/factory/contrib/afex/hystrix-go/v0"
	"github.com/xgodev/boost/middleware"
	"github.com/xgodev/boost/middleware/plugins/contrib/afex/hystrix-go/v0"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

func NewAnyError[R any](ctx context.Context, name string) middleware.AnyErrorMiddleware[R] {
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

func NewAny[R any](ctx context.Context, name string) middleware.AnyMiddleware[R] {
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

func NewError(ctx context.Context, name string) middleware.ErrorMiddleware {
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
