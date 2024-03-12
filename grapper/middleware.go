package grapper

import "context"

type ErrorExecFunc func(context.Context) error
type ErrorReturnFunc func(context.Context, error) error

type AnyErrorExecFunc[R any] func(context.Context) (R, error)
type AnyErrorReturnFunc[R any] func(context.Context, R, error) (R, error)

type AnyExecFunc[R any] func(context.Context) R
type AnyReturnFunc[R any] func(context.Context, R) R

type AnyErrorMiddleware[R any] interface {
	Exec(ctx *AnyErrorContext[R], exec AnyErrorExecFunc[R], returnFunc AnyErrorReturnFunc[R]) (R, error)
}

type AnyMiddleware[R any] interface {
	Exec(ctx *AnyContext[R], exec AnyExecFunc[R], returnFunc AnyReturnFunc[R]) R
}

type ErrorMiddleware interface {
	Exec(ctx *ErrorContext, exec ErrorExecFunc, returnFunc ErrorReturnFunc) error
}
