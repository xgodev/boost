package logger

import (
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

type Logger[T any] struct {
	options    *Options
	baseLogger log.Logger
}

func NewLogger[T any]() (*Logger[T], error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewLoggerWithOptions[T](opts), nil
}

func NewLoggerWithOptions[T any](options *Options) *Logger[T] {
	var zero T
	return &Logger[T]{options: options, baseLogger: log.WithTypeOf(zero)}
}

func NewAnyErrorMiddleware[T any]() (middleware.AnyErrorMiddleware[T], error) {
	return NewLogger[T]()
}

func NewAnyErrorMiddlewareWithOptions[T any](options *Options) middleware.AnyErrorMiddleware[T] {
	return NewLoggerWithOptions[T](options)
}

func (c *Logger[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (T, error) {
	logger := c.baseLogger.FromContext(ctx.GetContext())
	lm := c.logger(logger)

	e, err := ctx.Next(exec, fallbackFunc)
	if err != nil {
		logger.Error(err.Error())
		if c.options.ErrorStack {
			fmt.Println(errors.ErrorStack(err))
		}
	}

	var events []*event.Event

	switch r := any(e).(type) {
	case []*event.Event:
		events = r
	case *event.Event:
		if r == nil {
			return e, err
		}
		events = []*event.Event{r}
	default:
		return e, err
	}

	for _, ev := range events {
		j, err := json.Marshal(ev)
		if err != nil {
			logger.Errorf("error on marshall event for logging. %s", err.Error())
		} else {
			lm(string(j))
		}
	}

	return e, err
}

func (c *Logger[T]) logger(logger log.Logger) func(string) {
	switch c.options.Level {
	case "TRACE":
		return func(s string) { logger.Trace(s) }
	case "DEBUG":
		return func(s string) { logger.Debug(s) }
	default:
		return func(s string) { logger.Info(s) }
	}
}
