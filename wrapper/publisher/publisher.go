package publisher

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/extra/middleware"
)

type Publisher struct {
	driver Driver
	wrp    *middleware.AnyErrorWrapper[[]PublishOutput]
}

func New(driver Driver, mid ...middleware.AnyErrorMiddleware[[]PublishOutput]) *Publisher {
	wrp := middleware.NewAnyErrorWrapper[[]PublishOutput](context.Background(), "bootstrap", mid...)
	return &Publisher{driver: driver, wrp: wrp}
}

func (p *Publisher) Publish(ctx context.Context, events []*cloudevents.Event) error {
	_, err := p.wrp.Exec(ctx, "publisher",
		func(ctx context.Context) ([]PublishOutput, error) {
			return p.driver.Publish(ctx, events)
		}, nil)
	return err
}
