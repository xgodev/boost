package ignore_errors

import (
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"reflect"
)

type IgnoreErrors[T any] struct {
	options *Options
}

func (c *IgnoreErrors[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (T, error) {

	logger := log.FromContext(ctx.GetContext())

	e, err := ctx.Next(exec, fallbackFunc)
	if err != nil {

		err = errors.Cause(err)

		errType := reflect.TypeOf(err).Elem().Name()

		logger.Warnf("contains error type %s.  %s", errType, err.Error())

		for _, allowedErrorType := range c.options.Errors {

			if errType == allowedErrorType {
				logger.Warnf("ignoring error type %s.  %s", errType, err.Error())
				return e, nil
			}
		}
	}
	return e, err
}

func NewAnyErrorMiddleware[T any]() (middleware.AnyErrorMiddleware[T], error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewAnyErrorMiddlewareWithOptions[T](opts), nil
}

func NewAnyErrorMiddlewareWithOptions[T any](options *Options) middleware.AnyErrorMiddleware[T] {
	return &IgnoreErrors[T]{options: options}
}
