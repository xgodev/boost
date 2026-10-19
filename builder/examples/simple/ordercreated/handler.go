// Package ordercreated records the orders of order.created events in app
// worker.
package ordercreated

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/builder/examples/simple/order"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// RecordedType is the type of the event Handle replies with.
const RecordedType = "order.recorded"

// Handler records orders.
type Handler struct {
	svc *order.Service
}

// New provides the handler; the generated function door gets it from here.
// @Provide(index=0)
func New(svc *order.Service) *Handler {
	return &Handler{svc: svc}
}

type payload struct {
	Item string `json:"item"`
}

// Handle records the order of the event and replies with it.
// @Function(adapter=[cloudevents, nats], app=worker, middlewares=[recovery, logger])
func (h *Handler) Handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {
	var p payload
	if err := in.DataAs(&p); err != nil || p.Item == "" {
		return nil, errors.NotValidf("event %s has no item", in.ID())
	}

	o := h.svc.Create(p.Item)
	log.FromContext(ctx).Infof("recorded order %s of %s", o.ID, o.Item)

	out := cloudevents.NewEvent()
	out.SetType(RecordedType)
	out.SetSource("ordercreated")
	if err := out.SetData(cloudevents.ApplicationJSON, o); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error encoding order %s", o.ID))
	}
	return &out, nil
}
