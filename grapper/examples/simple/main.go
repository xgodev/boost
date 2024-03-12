package main

import (
	"context"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/xgodev/boost/grapper"
	h "github.com/xgodev/boost/grapper/middleware/contrib/afex/hystrix-go.v0"
	"github.com/xgodev/boost/grapper/middleware/contrib/americanas-go/log.v1"
	"github.com/xgodev/boost/log/contrib/sirupsen/logrus.v1"
)

func main() {

	ctx := context.Background()

	logrus.NewLogger()

	var r string
	var err error

	middlewares := []grapper.AnyErrorMiddleware[string]{
		log.NewAnyErrorMiddleware[string](ctx),
		h.NewAnyErrorMiddlewareWithConfig[string](ctx, "XPTO", hystrix.CommandConfig{
			Timeout:                10,
			MaxConcurrentRequests:  6000,
			RequestVolumeThreshold: 6000,
			SleepWindow:            10,
			ErrorPercentThreshold:  2,
		}),
	}

	wrp := grapper.NewAnyErrorWrapper[string](ctx, "example", middlewares...)

	r, err = wrp.Exec(ctx, "1", MyFunc, func(ctx context.Context, s string, err error) (string, error) {
		//send mail
		return "", err
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(r)
}

func MyFunc(ctx context.Context) (string, error) {
	return "string", nil
}
