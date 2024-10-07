package ignore_errors

import (
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"reflect"
	"strings"
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

		logger.Debugf("configured ignored error types: [%s]", strings.Join(c.options.Errors, ", "))
		logger.Warnf("contains error type %s. %s",
			errType,
			err.Error())

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
	return NewIgnoreErrors[T]()
}

func NewAnyErrorMiddlewareWithOptions[T any](options *Options) middleware.AnyErrorMiddleware[T] {
	return NewIgnoreErrorsWithOptions[T](options)
}

func NewIgnoreErrors[T any]() (*IgnoreErrors[T], error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewIgnoreErrorsWithOptions[T](opts), nil
}

func NewIgnoreErrorsWithOptions[T any](options *Options) *IgnoreErrors[T] {
	return &IgnoreErrors[T]{options: options}
}
