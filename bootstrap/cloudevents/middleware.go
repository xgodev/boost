package cloudevents

import (
	"context"

	v2 "github.com/cloudevents/sdk-go/v2"
)

// Middleware defines an interface to process middleware.
type Middleware interface {
	BeforeAll(ctx context.Context, inout []*InOut) (context.Context, error)
	Before(ctx context.Context, in *v2.Event) (context.Context, error)
	After(ctx context.Context, in v2.Event, out *v2.Event, err error) (context.Context, error)
	AfterAll(ctx context.Context, inout []*InOut) (context.Context, error)
	Close(ctx context.Context) error
}

// UnimplementedMiddleware defines default implementation for middleware.
type UnimplementedMiddleware struct {
}

// BeforeAll is called before all handlers, but this implementation does nothing.
func (u UnimplementedMiddleware) BeforeAll(ctx context.Context, inout []*InOut) (context.Context, error) {
	return ctx, nil
}

// Before is called before each handler, but this implementation does nothing.
func (u UnimplementedMiddleware) Before(ctx context.Context, in *v2.Event) (context.Context, error) {
	return ctx, nil
}

// After is called after each handler, but this implementation does nothing.
func (u UnimplementedMiddleware) After(ctx context.Context, in v2.Event, out *v2.Event, err error) (context.Context, error) {
	return ctx, nil
}

// AfterAll is called after all handlers, but this implementation does nothing.
func (u UnimplementedMiddleware) AfterAll(ctx context.Context, inout []*InOut) (context.Context, error) {
	return ctx, nil
}

// Close is called after the AfterAll hook to free any resources. In this implementation nothing is done.
func (u UnimplementedMiddleware) Close(ctx context.Context) error {
	return nil
}

// NewUnimplementedMiddleware returns a default middleware implementation.
func NewUnimplementedMiddleware() Middleware {
	return &UnimplementedMiddleware{}
}
