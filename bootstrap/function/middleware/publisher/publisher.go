package publisher

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/middleware"
	"github.com/xgodev/boost/wrapper/log"
)

type Publisher struct {
	driver Driver
}

func (c *Publisher) Exec(ctx *middleware.AnyErrorContext[*event.Event], exec middleware.AnyErrorExecFunc[*event.Event], fallbackFunc middleware.AnyErrorReturnFunc[*event.Event]) (*event.Event, error) {
	if !IsEnabled() {
		return ctx.Next(exec, fallbackFunc)
	}

	log.Tracef("publishing event")
	e, err := ctx.Next(exec, fallbackFunc)
	if err == nil {
		err = c.driver.Publish(ctx.GetContext(), []*event.Event{e})
		if err != nil {
			return nil, err
		}
	}
	return e, err
}

func New(driver Driver) middleware.AnyErrorMiddleware[*event.Event] {
	return &Publisher{driver: driver}
}
