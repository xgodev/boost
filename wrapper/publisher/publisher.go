package publisher

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/wrapper/log"
)

type Publisher struct {
	driver Driver
}

func New(driver Driver) *Publisher {
	return &Publisher{driver: driver}
}

func (p *Publisher) Publish(ctx context.Context, events []*cloudevents.Event) error {
	logger := log.FromContext(ctx).WithTypeOf(*p)
	logger.Tracef("publishing event")
	return p.driver.Publish(ctx, events)
}
