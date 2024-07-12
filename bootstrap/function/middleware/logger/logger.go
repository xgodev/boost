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

func (c *Logger) Exec(ctx *middleware.AnyErrorContext[*event.Event], exec middleware.AnyErrorExecFunc[*event.Event], fallbackFunc middleware.AnyErrorReturnFunc[*event.Event]) (*event.Event, error) {
	logger := log.FromContext(ctx.GetContext()).WithTypeOf(*c)
	lm := c.logger(logger)

	e, err := ctx.Next(exec, fallbackFunc)
	if err != nil {
		logger.Errorf(errors.ErrorStack(err))
		return e, err
	}

	if e == nil {
		j, err := json.Marshal(e)
		if err != nil {
			logger.Error(errors.ErrorStack(err))
		} else {
			lm(string(j))
		}
	}

	return e, err
}

func New(options *Options) middleware.AnyErrorMiddleware[*event.Event] {
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
