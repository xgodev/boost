package logger

import (
	"encoding/json"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

type Logger struct {
	options *Options
}

func (c *Logger) Exec(ctx *middleware.AnyErrorContext[any], exec middleware.AnyErrorExecFunc[any], fallbackFunc middleware.AnyErrorReturnFunc[any]) (any, error) {
	logger := log.FromContext(ctx.GetContext()).WithTypeOf(*c)
	lm := c.logger(logger)

	e, err := ctx.Next(exec, fallbackFunc)
	if err != nil {
		logger.Errorf(errors.ErrorStack(err))
		return e, err
	}

	if e != nil {

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
			j, err := json.Marshal(ev)
			if err != nil {
				logger.Error(errors.ErrorStack(err))
			} else {
				lm(string(j))
			}
		}

	}

	return e, err
}

func New(options *Options) middleware.AnyErrorMiddleware[any] {
	return &Logger{options: options}
}

func (c *Logger) logger(logger log.Logger) func(format string, args ...interface{}) {

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
