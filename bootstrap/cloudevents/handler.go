package cloudevents

import (
	"context"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Handler can be used to process events.
type Handler struct {
	handler *HandlerWrapper
}

// NewHandler creates a new handler wrapped in middleware.
func NewHandler(h *HandlerWrapper) *Handler {
	return &Handler{handler: h}
}

// Handle processes an event by calling the necessary middlewares.
func (h *Handler) Handle(ctx context.Context, in v2.Event) (out *v2.Event, err error) {

	logger := log.FromContext(ctx).WithTypeOf(*h)

	if isSQSEvent(ctx, &in) {
		fromSQS(ctx, &in)
	}

	inouts := []*InOut{
		{
			In:  &in,
			Err: err,
		},
	}

	err = h.handler.Process(ctx, inouts)
	if err != nil {
		logger.Error(errors.ErrorStack(err))
		return nil, err
	}

	logger.Debug("all events called")

	for _, inout := range inouts {
		if inout.Err != nil {
			err := errors.Wrap(inout.Err, errors.New("closing with errors lambda handle"))
			logger.Error(errors.ErrorStack(err))
			return nil, err
		}
	}

	return inouts[0].Out, nil
}
