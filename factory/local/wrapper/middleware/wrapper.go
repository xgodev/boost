package middleware

import (
	"context"
	"github.com/xgodev/boost/extra/middleware"
)

func NewAnyErrorWrapper[R any](ctx context.Context, plugins ...func(ctx context.Context, name string) middleware.AnyErrorMiddleware[R]) (*middleware.AnyErrorWrapper[R], error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewAnyErrorWrapperWithOptions(ctx, o, plugins...)
}

func NewAnyErrorWrapperWithPath[R any](ctx context.Context, path string, plugins ...func(ctx context.Context, name string) middleware.AnyErrorMiddleware[R]) (*middleware.AnyErrorWrapper[R], error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewAnyErrorWrapperWithOptions(ctx, o, plugins...)
}

func NewAnyErrorWrapperWithOptions[R any](ctx context.Context, options *Options, plugins ...func(ctx context.Context, name string) middleware.AnyErrorMiddleware[R]) (*middleware.AnyErrorWrapper[R], error) {
	var m []middleware.AnyErrorMiddleware[R]

	for _, f := range plugins {
		if p := f(ctx, options.Name); p != nil {
			m = append(m, p)
		}
	}

	return middleware.NewAnyErrorWrapper(ctx, options.Name, m...), nil
}

func NewAnyWrapper[R any](ctx context.Context, plugins ...func(ctx context.Context, name string) middleware.AnyMiddleware[R]) (*middleware.AnyWrapper[R], error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewAnyWrapperWithOptions(ctx, o, plugins...)
}

func NewAnyWrapperWithPath[R any](ctx context.Context, path string, plugins ...func(ctx context.Context, name string) middleware.AnyMiddleware[R]) (*middleware.AnyWrapper[R], error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewAnyWrapperWithOptions(ctx, o, plugins...)
}

func NewAnyWrapperWithOptions[R any](ctx context.Context, options *Options, plugins ...func(ctx context.Context, name string) middleware.AnyMiddleware[R]) (*middleware.AnyWrapper[R], error) {
	var m []middleware.AnyMiddleware[R]

	for _, f := range plugins {
		if p := f(ctx, options.Name); p != nil {
			m = append(m, p)
		}
	}

	return middleware.NewAnyWrapper(ctx, options.Name, m...), nil
}

func NewErrorWrapper(ctx context.Context, plugins ...func(ctx context.Context, name string) middleware.ErrorMiddleware) (*middleware.ErrorWrapper, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewErrorWrapperWithOptions(ctx, o, plugins...)
}

func NewErrorWrapperWithPath(ctx context.Context, path string, plugins ...func(ctx context.Context, name string) middleware.ErrorMiddleware) (*middleware.ErrorWrapper, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewErrorWrapperWithOptions(ctx, o, plugins...)
}

func NewErrorWrapperWithOptions(ctx context.Context, options *Options, plugins ...func(ctx context.Context, name string) middleware.ErrorMiddleware) (*middleware.ErrorWrapper, error) {
	var m []middleware.ErrorMiddleware

	for _, f := range plugins {
		if p := f(ctx, options.Name); p != nil {
			m = append(m, p)
		}
	}

	return middleware.NewErrorWrapper(ctx, options.Name, m...), nil
}
