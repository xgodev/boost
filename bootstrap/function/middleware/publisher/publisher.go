package publisher

import (
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"

	stderrors "errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
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
		if r == nil {
			return e, err
		}
		events = []*event.Event{r}
	case nil:
		return e, err
	default:
		return e, errors.Internalf("unsupported handler type")
	}

	var deadLetterSubject string
	var errorType string

	if err != nil && c.options.Deadletter.Enabled {

		logger.Debugf("configured to send to deadletter error types: [%s]", strings.Join(c.options.Deadletter.Errors, ", "))

		if ok, name := shouldIgnoreError(err, c.options.Deadletter.Errors); ok {
			logger.Warnf("Error type %s is allowed to be sent to dead letter", name)
			deadLetterSubject = c.options.Deadletter.Subject
			errorType = name
		} else {
			return e, err
		}

		if deadLetterSubject == "" {
			logger.Warnf("no dead letter subject found for error type %s", errorType)
			return e, err
		}

		for _, ev := range events {
			ev.SetSubject(deadLetterSubject)
			ev.SetExtension("error_type", errorType)
			ev.SetExtension("error", err.Error())
		}

	} else {
		logger.Debugf("dead letter is disabled. ignoring dead letter")
		return e, err
	}

	for _, ev := range events {
		if ev.Subject() == "" {
			if c.options.Subject == "" {
				logger.Warnf("no subject found for event. ignoring publish")
				return e, err
			}
			ev.SetSubject(c.options.Subject)
		}
	}

	if len(events) == 0 {
		logger.Debugf("no events to publish")
		return e, err
	}

	if err := c.publisher.Publish(ctx.GetContext(), events); err != nil {
		logger.Errorf("error publishing event: %v", err)
		return e, err
	}

	return e, err
}

func shouldIgnoreError(err error, allowed []string) (bool, string) {
	for err != nil {
		errName := fmt.Sprintf("%T", err)          // ex: *my.ErrFoo
		errName = strings.TrimPrefix(errName, "*") // remove o '*' para comparar com o nome puro

		for _, allowedName := range allowed {
			if strings.HasSuffix(errName, allowedName) {
				return true, errName
			}
		}

		err = stderrors.Unwrap(err)
	}
	return false, ""
}

func NewAnyErrorMiddleware[T any](publisher *publisher.Publisher) (middleware.AnyErrorMiddleware[T], error) {
	return NewPublisher[T](publisher)
}

func NewAnyErrorMiddlewareWithOptions[T any](publisher *publisher.Publisher, options *Options) middleware.AnyErrorMiddleware[T] {
	return NewPublisherWithOptions[T](publisher, options)
}

func NewPublisher[T any](publisher *publisher.Publisher) (*Publisher[T], error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewPublisherWithOptions[T](publisher, opts), nil
}

func NewPublisherWithOptions[T any](publisher *publisher.Publisher, options *Options) *Publisher[T] {
	return &Publisher[T]{publisher: publisher, options: options}
}
