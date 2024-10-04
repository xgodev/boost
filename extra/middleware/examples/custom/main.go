package main

import (
	"context"
	"fmt"
	"github.com/xgodev/boost/extra/middleware"
)

type CustomMiddleware[R any] struct{}

func (c *CustomMiddleware[R]) Exec(ctx *middleware.AnyErrorContext[R], exec middleware.AnyErrorExecFunc[R], fallbackFunc middleware.AnyErrorReturnFunc[R]) (R, error) {
	fmt.Println("my custom middleware")
	return ctx.Next(exec, fallbackFunc)
}

func NewCustomMiddleware[R any]() middleware.AnyErrorMiddleware[R] {
	return &CustomMiddleware[R]{}
}

func main() {

	ctx := context.Background()

	var res string
	var err error

	middlewares := []middleware.AnyErrorMiddleware[string]{
		NewCustomMiddleware[string](),
	}

	wrp := middleware.NewAnyErrorWrapper[string](ctx, "example", middlewares...)

	res, err = wrp.Exec(ctx, "1",
		func(ctx context.Context) (string, error) {
			return "string", nil
		}, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res)
}
