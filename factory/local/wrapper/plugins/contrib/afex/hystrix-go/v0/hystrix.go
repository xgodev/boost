package hystrix

import (
	"context"
	"github.com/xgodev/boost/config"
	h "github.com/xgodev/boost/factory/contrib/afex/hystrix-go/v0"
	"github.com/xgodev/boost/log"
	"github.com/xgodev/boost/wrapper"
	"github.com/xgodev/boost/wrapper/middleware/contrib/afex/hystrix-go/v0"
)

func NewAnyError[R any](ctx context.Context, name string) wrapper.AnyErrorMiddleware[R] {
	ConfigAdd(name)
	config.Load()
	if o, _ := NewOptions(name); !o.Enabled {
		return nil
	}
	if err := h.ConfigureCommand(name); err != nil {
		log.Error(err.Error())
	}
	return hystrix.hystrix.NewAnyErrorMiddleware[R](ctx, name)
}

func NewAny[R any](ctx context.Context, name string) wrapper.AnyMiddleware[R] {
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

func NewError(ctx context.Context, name string) wrapper.ErrorMiddleware {
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
