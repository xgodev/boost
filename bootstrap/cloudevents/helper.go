package cloudevents

import (
	"context"

	cloudevents "github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2"
)

// Helper assists in creating event handlers.
type Helper struct {
	ctx    context.Context
	client *cloudevents.Client
}

// NewHelper returns a new Helper.
func NewHelper(ctx context.Context, handler *HandlerWrapper) *Helper {

	client := cloudevents.NewHTTP(ctx, NewHandler(handler).Handle)

	return &Helper{
		client: client,
		ctx:    ctx,
	}
}

// Start starts HTTP client for handle events.
func (h *Helper) Start() {
	h.client.Start(h.ctx)
}
