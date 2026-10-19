// Package ordercreated records the orders of order.created events in app
// worker.
package ordercreated

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/builder/examples/simple/order"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
)

// RecordedType is the type of the event Handle replies with.
const RecordedType = "order.recorded"

// Handler records orders.
type Handler struct {
	svc *order.Service
	pub *publisher.Publisher
}

// New provides the handler; the generated function door gets it from here.
// The publisher announces the recorded orders, on the broker the function
// may read from; only the nats and pubsub drivers (and noop, the default)
// are generated, so the worker links no Kafka client nor librdkafka.
// @Provide(index=0)
// @Inject(index=1, name=orders, drivers=[nats, pubsub])
func New(svc *order.Service, pub *publisher.Publisher) *Handler {
	return &Handler{svc: svc, pub: pub}
}

type payload struct {
	Item string `json:"item"`
}

// Handle records the order of the event, announces it and replies with it.
// @Function(adapter=[cloudevents, nats, pubsub], app=worker, client=orders, middlewares=[recovery, logger])
func (h *Handler) Handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {
	var p payload
	if err := in.DataAs(&p); err != nil || p.Item == "" {
		return nil, errors.NotValidf("event %s has no item", in.ID())
	}

	o := h.svc.Create(p.Item)
	log.FromContext(ctx).Infof("recorded order %s of %s", o.ID, o.Item)

	out := cloudevents.NewEvent()
	out.SetID(o.ID)
	out.SetType(RecordedType)
	out.SetSource("ordercreated")
	if err := out.SetData(cloudevents.ApplicationJSON, o); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error encoding order %s", o.ID))
	}
	if err := h.pub.Publish(ctx, []*cloudevents.Event{&out}); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error announcing order %s", o.ID))
	}
	return &out, nil
}
