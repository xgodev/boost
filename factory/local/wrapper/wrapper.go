package wrapper

import (
	"context"
	"github.com/xgodev/boost/wrapper"
)

func NewAnyErrorWrapper[R any](ctx context.Context, plugins ...func(ctx context.Context, name string) wrapper.AnyErrorMiddleware[R]) (*wrapper.AnyErrorWrapper[R], error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewAnyErrorWrapperWithOptions(ctx, o, plugins...)
}

func NewAnyErrorWrapperWithPath[R any](ctx context.Context, path string, plugins ...func(ctx context.Context, name string) wrapper.AnyErrorMiddleware[R]) (*wrapper.AnyErrorWrapper[R], error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewAnyErrorWrapperWithOptions(ctx, o, plugins...)
}

func NewAnyErrorWrapperWithOptions[R any](ctx context.Context, options *Options, plugins ...func(ctx context.Context, name string) wrapper.AnyErrorMiddleware[R]) (*wrapper.AnyErrorWrapper[R], error) {
	var m []wrapper.AnyErrorMiddleware[R]

	for _, f := range plugins {
		if p := f(ctx, options.Name); p != nil {
			m = append(m, p)
		}
	}

	return wrapper.NewAnyErrorWrapper(ctx, options.Name, m...), nil
}

func NewAnyWrapper[R any](ctx context.Context, plugins ...func(ctx context.Context, name string) wrapper.AnyMiddleware[R]) (*wrapper.AnyWrapper[R], error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewAnyWrapperWithOptions(ctx, o, plugins...)
}

func NewAnyWrapperWithPath[R any](ctx context.Context, path string, plugins ...func(ctx context.Context, name string) wrapper.AnyMiddleware[R]) (*wrapper.AnyWrapper[R], error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewAnyWrapperWithOptions(ctx, o, plugins...)
}

func NewAnyWrapperWithOptions[R any](ctx context.Context, options *Options, plugins ...func(ctx context.Context, name string) wrapper.AnyMiddleware[R]) (*wrapper.AnyWrapper[R], error) {
	var m []wrapper.AnyMiddleware[R]

	for _, f := range plugins {
		if p := f(ctx, options.Name); p != nil {
			m = append(m, p)
		}
	}

	return wrapper.NewAnyWrapper(ctx, options.Name, m...), nil
}

func NewErrorWrapper(ctx context.Context, plugins ...func(ctx context.Context, name string) wrapper.ErrorMiddleware) (*wrapper.ErrorWrapper, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewErrorWrapperWithOptions(ctx, o, plugins...)
}

func NewErrorWrapperWithPath(ctx context.Context, path string, plugins ...func(ctx context.Context, name string) wrapper.ErrorMiddleware) (*wrapper.ErrorWrapper, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewErrorWrapperWithOptions(ctx, o, plugins...)
}

func NewErrorWrapperWithOptions(ctx context.Context, options *Options, plugins ...func(ctx context.Context, name string) wrapper.ErrorMiddleware) (*wrapper.ErrorWrapper, error) {
	var m []wrapper.ErrorMiddleware

	for _, f := range plugins {
		if p := f(ctx, options.Name); p != nil {
			m = append(m, p)
		}
	}

	return wrapper.NewErrorWrapper(ctx, options.Name, m...), nil
}
