package publisher

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/publisher"
)

type Publisher struct {
	publisher *publisher.Publisher
	options   *Options
}

func (c *Publisher) Exec(ctx *middleware.AnyErrorContext[any], exec middleware.AnyErrorExecFunc[any], fallbackFunc middleware.AnyErrorReturnFunc[any]) (any, error) {
	e, err := ctx.Next(exec, fallbackFunc)
	if err == nil && e != nil {

		var events []*event.Event

		switch r := e.(type) {
		case []*event.Event:
			events = r
		case *event.Event:
			events = []*event.Event{r}
		default:
			return nil, errors.Errorf("unsupported handler type")
		}

		for _, ev := range events {
			if ev.Subject() == "" {
				ev.SetSubject(c.options.Subject)
			}
		}

		err = c.publisher.Publish(ctx.GetContext(), events)
		if err != nil {
			return nil, err
		}
	}
	return e, err
}

func New(publisher *publisher.Publisher, options *Options) middleware.AnyErrorMiddleware[any] {
	return &Publisher{publisher: publisher, options: options}
}
