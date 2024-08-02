package function

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
)

func Wrapper(wrp *middleware.AnyErrorWrapper[any], fn Handler) Handler {
	return func(ctx context.Context, in event.Event) (any, error) {
		return wrp.Exec(ctx, Name(),
			func(ctx context.Context) (any, error) {
				return fn(ctx, in)
			}, nil)
	}
}
