package recovery

import (
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

type Recover[T any] struct {
}

func NewRecover[T any]() *Recover[T] {
	return &Recover[T]{}
}

func NewAnyErrorMiddleware[T any]() middleware.AnyErrorMiddleware[T] {
	return NewRecover[T]()
}

func (c *Recover[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (res T, err error) {
	logger := log.FromContext(ctx.GetContext()).WithTypeOf(*c)

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("recovering %v", r)
			err = errors.Internalf("%v", r)
		}
	}()

	res, err = ctx.Next(exec, fallbackFunc)
	return res, err
}
