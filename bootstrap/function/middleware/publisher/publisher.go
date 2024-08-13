package publisher

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/wrapper/publisher"
)

type Publisher[T any] struct {
	publisher *publisher.Publisher
	options   *Options
}

func (c *Publisher[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (T, error) {
	e, err := ctx.Next(exec, fallbackFunc)

	var events []*event.Event

	switch r := any(e).(type) {
	case []*event.Event:
		events = r
	case *event.Event:
		events = []*event.Event{r}
	default:
		return e, err
	}

	for _, ev := range events {
		if ev == nil {
			continue
		}
		if ev.Subject() == "" {
			ev.SetSubject(c.options.Subject)
		}
	}

	return e, c.publisher.Publish(ctx.GetContext(), events)

}

func NewAnyErrorMiddleware[T any](publisher *publisher.Publisher) (middleware.AnyErrorMiddleware[T], error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewAnyErrorMiddlewareWithOptions[T](publisher, opts), nil
}

func NewAnyErrorMiddlewareWithOptions[T any](publisher *publisher.Publisher, options *Options) middleware.AnyErrorMiddleware[T] {
	return &Publisher[T]{publisher: publisher, options: options}
}
