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
	options *Options
}

func NewLogger[T any]() (*Logger[T], error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewLoggerWithOptions[T](opts), nil
}

func NewLoggerWithOptions[T any](options *Options) *Logger[T] {
	return &Logger[T]{options: options}
}

func NewAnyErrorMiddleware[T any]() (middleware.AnyErrorMiddleware[T], error) {
	return NewLogger[T]()
}

func NewAnyErrorMiddlewareWithOptions[T any](options *Options) middleware.AnyErrorMiddleware[T] {
	return NewLoggerWithOptions[T](options)
}

func (c *Logger[T]) Exec(ctx *middleware.AnyErrorContext[T], exec middleware.AnyErrorExecFunc[T], fallbackFunc middleware.AnyErrorReturnFunc[T]) (T, error) {
	logger := log.FromContext(ctx.GetContext()).WithTypeOf(*c)
	lm := c.logger(logger)

	e, err := ctx.Next(exec, fallbackFunc)
	var events []*event.Event

	switch r := any(e).(type) {
	case []*event.Event:
		events = r
	case *event.Event:
		if r != nil {
			events = []*event.Event{r}
		}
		events = []*event.Event{r}
	default:
		return e, err
	}

	for _, ev := range events {
		if ev == nil {
			continue
		}
		j, err := json.Marshal(ev)
		if err != nil {
			logger.Error(err.Error())
			if c.options.ErrorStack {
				fmt.Println(errors.ErrorStack(err))
			}
		} else {
			lm(string(j))
		}
	}

	return e, err
}

func (c *Logger[T]) logger(logger log.Logger) func(format string, args ...interface{}) {

	var method func(format string, args ...interface{})

	switch c.options.Level {
	case "TRACE":
		method = logger.Tracef
	case "DEBUG":
		method = logger.Debugf
	default:
		method = logger.Infof
	}

	return method
}
