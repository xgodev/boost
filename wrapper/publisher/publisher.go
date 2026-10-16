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

// Publish publishes events and returns an error when the driver fails or when
// any event fails. Per-event failures are turned into the error only after the
// middleware chain, so middlewares still see them as outputs.
func (p *Publisher) Publish(ctx context.Context, events []*cloudevents.Event) error {
	outputs, err := p.wrp.Exec(ctx, "publisher",
		func(ctx context.Context) ([]PublishOutput, error) {
			return p.driver.Publish(ctx, events)
		}, nil)
	if err != nil {
		return err
	}
	return outputsError(outputs)
}
