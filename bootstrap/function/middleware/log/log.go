package log

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

type Log struct {
	options *Options
}

func (c *Log) Exec(ctx *middleware.AnyErrorContext[*event.Event], exec middleware.AnyErrorExecFunc[*event.Event], fallbackFunc middleware.AnyErrorReturnFunc[*event.Event]) (*event.Event, error) {
	logger := log.FromContext(ctx.GetContext()).WithTypeOf(*c)
	e, err := ctx.Next(exec, fallbackFunc)
	if err != nil {
		logger.Errorf(errors.ErrorStack(err))
	}
	return e, err
}

func New(options *Options) middleware.AnyErrorMiddleware[*event.Event] {
	return &Log{options: options}
}
