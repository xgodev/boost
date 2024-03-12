package grapper

import (
	"context"
)

type AnyErrorWrapper[R any] struct {
	m    []AnyErrorMiddleware[R]
	name string
}

func NewAnyErrorWrapper[R any](ctx context.Context, name string, m ...AnyErrorMiddleware[R]) *AnyErrorWrapper[R] {
	return &AnyErrorWrapper[R]{m: m, name: name}
}

func (w *AnyErrorWrapper[R]) Exec(ctx context.Context, id string, exec AnyErrorExecFunc[R], returnFunc AnyErrorReturnFunc[R]) (r R, err error) {
	c := NewAnyErrorContext(w.name, id, w.m...)
	c.SetContext(ctx)
	return c.Next(exec, returnFunc)
}

type AnyWrapper[R any] struct {
	m    []AnyMiddleware[R]
	name string
}

func NewAnyWrapper[R any](ctx context.Context, name string, m ...AnyMiddleware[R]) *AnyWrapper[R] {
	return &AnyWrapper[R]{m: m, name: name}
}

func (w *AnyWrapper[R]) Exec(ctx context.Context, id string, exec AnyExecFunc[R], returnFunc AnyReturnFunc[R]) (r R) {
	c := NewAnyContext(w.name, id, w.m...)
	c.SetContext(ctx)
	return c.Next(exec, returnFunc)
}

type ErrorWrapper struct {
	m    []ErrorMiddleware
	name string
}

func NewErrorWrapper(ctx context.Context, name string, m ...ErrorMiddleware) *ErrorWrapper {
	return &ErrorWrapper{m: m, name: name}
}

func (w *ErrorWrapper) Exec(ctx context.Context, id string, exec ErrorExecFunc, returnFunc ErrorReturnFunc) (err error) {
	c := NewErrorContext(w.name, id, w.m...)
	c.SetContext(ctx)
	return c.Next(exec, returnFunc)
}
