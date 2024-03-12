package main

import (
	"context"
	"fmt"
	"github.com/xgodev/boost/log/contrib/sirupsen/logrus/v1"
	h "github.com/xgodev/boost/wrapper/middleware/contrib/afex/hystrix-go/v0"
	"github.com/xgodev/boost/wrapper/middleware/local/log"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/xgodev/boost/wrapper"
)

type Result struct {
	Code string
}

type FooService struct {
	wrapper *wrapper.AnyErrorWrapper[Result]
}

func NewFooService(wrapper *wrapper.AnyErrorWrapper[Result]) *FooService {
	return &FooService{wrapper: wrapper}
}

func (s *FooService) FooMethod(ctx context.Context) (Result, error) {
	return s.wrapper.Exec(ctx, "1", func(ctx context.Context) (Result, error) {
		return Result{Code: "XPTO"}, nil
	}, nil)
}

func main() {

	ctx := context.Background()

	logrus.NewLogger()

	var r Result
	var err error

	middlewares := []wrapper.AnyErrorMiddleware[Result]{
		log.NewAnyErrorMiddleware[Result](ctx),
		h.NewAnyErrorMiddlewareWithConfig[Result](ctx, "XPTO", hystrix.CommandConfig{
			Timeout:                10,
			MaxConcurrentRequests:  6000,
			RequestVolumeThreshold: 6000,
			SleepWindow:            10,
			ErrorPercentThreshold:  2,
		}),
	}

	wrapper := wrapper.NewAnyErrorWrapper[Result](ctx, "example", middlewares...)

	foo := NewFooService(wrapper)
	r, err = foo.FooMethod(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(r.Code)
}
