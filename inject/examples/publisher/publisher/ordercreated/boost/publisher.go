// Package boost implements ordercreated.Publisher on boost's
// wrapper/publisher, whose driver config chooses.
package boost

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/inject/examples/publisher/domain/model/entity"
	"github.com/xgodev/boost/inject/examples/publisher/domain/publisher/ordercreated"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/publisher"
)

const (
	// Subject is the NATS subject, Pub/Sub topic or Kafka topic of the event.
	Subject = "orders.created"
	// Type is the CloudEvents type of the event.
	Type = "order.created"
)

// payload is the event's data; the entity carries no JSON tags.
type payload struct {
	ID   string `json:"id"`
	Item string `json:"item"`
}

type eventPublisher struct {
	publisher *publisher.Publisher
}

// New returns the ordercreated.Publisher on the *publisher.Publisher inject
// builds with the driver named by boost.wrapper.publisher.driver.name.
// @Provide(index=0)
// @Inject(index=0)
func New(p *publisher.Publisher) ordercreated.Publisher {
	return &eventPublisher{publisher: p}
}

// Publish sends order as an order.created CloudEvent.
func (p *eventPublisher) Publish(ctx context.Context, order entity.Order) error {
	event := cloudevents.NewEvent()
	event.SetID(order.ID)
	event.SetType(Type)
	event.SetSource("inject/examples/publisher")
	event.SetSubject(Subject)
	if err := event.SetData(cloudevents.ApplicationJSON, payload{ID: order.ID, Item: order.Item}); err != nil {
		return errors.Wrap(err, errors.Internalf("error encoding order %s", order.ID))
	}
	return p.publisher.Publish(ctx, []*cloudevents.Event{&event})
}
