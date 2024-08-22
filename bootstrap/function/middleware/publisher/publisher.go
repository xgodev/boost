package publisher

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
	berrors "github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"reflect"
)

type Publisher[T any] struct {
	publisher *publisher.Publisher
	options   *Options
}

func (c *Publisher[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (T, error) {

	logger := log.FromContext(ctx.GetContext())

	e, err := ctx.Next(exec, fallbackFunc)
	if &e != nil {

		var events []*event.Event

		switch r := any(e).(type) {
		case []*event.Event:
			events = r
		case *event.Event:
			events = []*event.Event{r}
		default:
			return e, berrors.Internalf("unsupported handler type")
		}

		var deadLetterSubject string
		var errorType string

		if err != nil {

			if c.options.Deadletter.Enabled {

				errType := reflect.TypeOf(err).Elem().Name()

				for _, allowedErrorType := range c.options.Deadletter.Errors {
					if errType == allowedErrorType {
						logger.Tracef("Error type %s is allowed to be sent to dead letter", errType)
						deadLetterSubject = c.options.Deadletter.Subject
						errorType = errType
						break
					}
				}

			}

			if deadLetterSubject == "" {
				return e, err
			}

		}

		for _, ev := range events {

			if deadLetterSubject != "" {
				ev.SetSubject(deadLetterSubject)
				ev.SetExtension("error_type", errorType)
				ev.SetExtension("error", err.Error())
			} else if ev.Subject() == "" {
				ev.SetSubject(c.options.Subject)
			}
		}

		return e, c.publisher.Publish(ctx.GetContext(), events)
	}
	return e, err
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
