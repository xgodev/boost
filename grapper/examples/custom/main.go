package main

import (
	"context"
	"fmt"

	"github.com/xgodev/boost/grapper"
	"github.com/xgodev/boost/log/contrib/sirupsen/logrus.v1"
)

type CustomMiddleware[R any] struct{}

func (c *CustomMiddleware[R]) Exec(ctx *grapper.AnyErrorContext[R], exec grapper.AnyErrorExecFunc[R], fallbackFunc grapper.AnyErrorReturnFunc[R]) (R, error) {
	fmt.Println("my custom middleware")
	return ctx.Next(exec, fallbackFunc)
}

func NewCustomMiddleware[R any]() grapper.AnyErrorMiddleware[R] {
	return &CustomMiddleware[R]{}
}

func main() {

	ctx := context.Background()

	logrus.NewLogger()

	var res string
	var err error

	middlewares := []grapper.AnyErrorMiddleware[string]{
		NewCustomMiddleware[string](),
	}

	wrp := grapper.NewAnyErrorWrapper[string](ctx, "example", middlewares...)

	res, err = wrp.Exec(ctx, "1",
		func(ctx context.Context) (string, error) {
			return "string", nil
		}, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res)
}
