package ignore_errors

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/wrapper/log"
	"strings"
)

type IgnoreErrors[T any] struct {
	options *Options
}

func (c *IgnoreErrors[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (T, error) {
	logger := log.FromContext(ctx.GetContext())

	e, err := ctx.Next(exec, fallbackFunc)
	if err != nil {
		logger.Debugf("configured ignored error types: [%s]", strings.Join(c.options.Errors, ", "))

		// Verifica se algum erro na cadeia deve ser ignorado
		if shouldIgnoreError(err, c.options.Errors) {
			logger.Warnf("ignoring error: %s", err.Error())
			return e, nil
		}

	}

	return e, err
}

func shouldIgnoreError(err error, allowed []string) bool {
	for err != nil {
		errName := fmt.Sprintf("%T", err)          // ex: *my.ErrFoo
		errName = strings.TrimPrefix(errName, "*") // remove o '*' para comparar com o nome puro

		for _, allowedName := range allowed {
			if strings.HasSuffix(errName, allowedName) {
				return true
			}
		}

		err = errors.Unwrap(err)
	}
	return false
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
