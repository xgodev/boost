package main

import (
	"context"
	"fmt"
	"github.com/xgodev/boost/log/contrib/sirupsen/logrus/v1"

	"github.com/xgodev/boost/wrapper"
)

type CustomMiddleware[R any] struct{}

func (c *CustomMiddleware[R]) Exec(ctx *wrapper.AnyErrorContext[R], exec wrapper.AnyErrorExecFunc[R], fallbackFunc wrapper.AnyErrorReturnFunc[R]) (R, error) {
	fmt.Println("my custom middleware")
	return ctx.Next(exec, fallbackFunc)
}

func NewCustomMiddleware[R any]() wrapper.AnyErrorMiddleware[R] {
	return &CustomMiddleware[R]{}
}

func main() {

	ctx := context.Background()

	logrus.NewLogger()

	var res string
	var err error

	middlewares := []wrapper.AnyErrorMiddleware[string]{
		NewCustomMiddleware[string](),
	}

	wrp := wrapper.NewAnyErrorWrapper[string](ctx, "example", middlewares...)

	res, err = wrp.Exec(ctx, "1",
		func(ctx context.Context) (string, error) {
			return "string", nil
		}, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res)
}
