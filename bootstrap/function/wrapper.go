package function

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/middleware"
)

func Wrapper(wrp *middleware.AnyErrorWrapper[*event.Event], fn Handler) Handler {
	return func(ctx context.Context, in event.Event) (*event.Event, error) {
		return wrp.Exec(ctx, "func",
			func(ctx context.Context) (*event.Event, error) {
				return fn(ctx, in)
			}, nil)
	}
}
