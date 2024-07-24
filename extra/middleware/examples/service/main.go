package main

import (
	"context"
	"fmt"
	h "github.com/xgodev/boost/extra/middleware/plugins/contrib/afex/hystrix-go/v0"
	l "github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus/v1"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/xgodev/boost/extra/middleware"
)

type Result struct {
	Code string
}

type FooService struct {
	wrapper *middleware.AnyErrorWrapper[Result]
}

func NewFooService(wrapper *middleware.AnyErrorWrapper[Result]) *FooService {
	return &FooService{wrapper: wrapper}
}

func (s *FooService) FooMethod(ctx context.Context) (Result, error) {
	return s.wrapper.Exec(ctx, "1", func(ctx context.Context) (Result, error) {
		return Result{Code: "XPTO"}, nil
	}, nil)
}

func main() {

	ctx := context.Background()

	log.Set(logrus.NewLogger())

	var r Result
	var err error

	middlewares := []middleware.AnyErrorMiddleware[Result]{
		l.NewAnyErrorMiddleware[Result](ctx),
		h.NewAnyErrorMiddlewareWithConfig[Result](ctx, "XPTO", hystrix.CommandConfig{
			Timeout:                10,
			MaxConcurrentRequests:  6000,
			RequestVolumeThreshold: 6000,
			SleepWindow:            10,
			ErrorPercentThreshold:  2,
		}),
	}

	wrapper := middleware.NewAnyErrorWrapper[Result](ctx, "example", middlewares...)

	foo := NewFooService(wrapper)
	r, err = foo.FooMethod(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(r.Code)
}
