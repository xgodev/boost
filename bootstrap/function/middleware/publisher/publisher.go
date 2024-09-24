package publisher

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"reflect"
	"strings"
)

type Publisher[T any] struct {
	publisher *publisher.Publisher
	options   *Options
}

func (c *Publisher[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (T, error) {

	logger := log.FromContext(ctx.GetContext())

	e, err := ctx.Next(exec, fallbackFunc)

	var events []*event.Event

	switch r := any(e).(type) {
	case []*event.Event:
		events = r
	case *event.Event:
		events = []*event.Event{r}
	case nil:
		events = []*event.Event{}
	default:
		return e, errors.Internalf("unsupported handler type")
	}

	var deadLetterSubject string
	var errorType string

	if err != nil {

		if c.options.Deadletter.Enabled {

			err = errors.Cause(err)

			errType := reflect.TypeOf(err).Elem().Name()

			logger.Debugf("configured to send to deadletter error types: [%s]", strings.Join(c.options.Deadletter.Errors, ", "))
			logger.Warnf("contains error type %s. %s",
				errType,
				err.Error())

			for _, allowedErrorType := range c.options.Deadletter.Errors {
				if errType == allowedErrorType {
					logger.Tracef("Error type %s is allowed to be sent to dead letter", errType)
					deadLetterSubject = c.options.Deadletter.Subject
					errorType = errType
					break
				}
			}

			if deadLetterSubject == "" {
				logger.Warnf("no dead letter subject found for error type %s", errType)
				return e, nil
			}

			for _, ev := range events {
				ev.SetSubject(deadLetterSubject)
				ev.SetExtension("error_type", errorType)
				ev.SetExtension("error", err.Error())
			}

		} else {
			logger.Debugf("dead letter is disabled. ignoring dead letter")
		}

	}

	for _, ev := range events {
		if ev.Subject() == "" {
			if c.options.Subject == "" {
				logger.Warnf("no subject found for event. ignoring publish")
				return e, nil
			}
			ev.SetSubject(c.options.Subject)
		}
	}
	
	if len(events) == 0 {
		logger.Debugf("no events to publish")
		return e, nil
	}

	if err := c.publisher.Publish(ctx.GetContext(), events); err != nil {
		return e, err
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
