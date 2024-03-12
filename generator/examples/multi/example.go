package multi

import (
	"context"

	v2 "github.com/cloudevents/sdk-go/v2"
	igce "github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2"
)

func ExampleHandler() igce.Handler {
	h := NewHandler()
	return h.Handle
}

type exampleHandler struct {
}

func (h *exampleHandler) Handle(ctx context.Context, in v2.Event) (out *v2.Event, err error) {
	return nil, err
}

func NewHandler() *exampleHandler {
	return &exampleHandler{}
}
