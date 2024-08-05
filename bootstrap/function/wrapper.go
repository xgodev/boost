package function

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
)

func Wrapper[T any](wrp *middleware.AnyErrorWrapper[T], fn Handler[T]) Handler[T] {
	return func(ctx context.Context, in event.Event) (T, error) {
		return wrp.Exec(ctx, Name(),
			func(ctx context.Context) (T, error) {
				return fn(ctx, in)
			}, nil)
	}
}
