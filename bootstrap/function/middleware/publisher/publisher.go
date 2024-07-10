package publisher

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/wrapper/publisher"
)

type Publisher struct {
	publisher *publisher.Publisher
	options   *Options
}

func (c *Publisher) Exec(ctx *middleware.AnyErrorContext[*event.Event], exec middleware.AnyErrorExecFunc[*event.Event], fallbackFunc middleware.AnyErrorReturnFunc[*event.Event]) (*event.Event, error) {
	e, err := ctx.Next(exec, fallbackFunc)
	if err == nil && e != nil {
		e.SetSubject(c.options.Subject)
		err = c.publisher.Publish(ctx.GetContext(), []*event.Event{e})
		if err != nil {
			return nil, err
		}
	}
	return e, err
}

func New(publisher *publisher.Publisher, options *Options) middleware.AnyErrorMiddleware[*event.Event] {
	return &Publisher{publisher: publisher, options: options}
}
