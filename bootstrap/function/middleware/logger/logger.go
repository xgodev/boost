package logger

import (
	"fmt"

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
	//logCtx := zerolog.NewLogger().ToContext(ctx.GetContext())
	//ctx.SetContext(logCtx)

	e, err := ctx.Next(exec, fallbackFunc)
	if err != nil {
		log.Ctx(ctx.GetContext(), *c).Warnf("handle with error: %s", err.Error())
		if c.options.ErrorStack {
			fmt.Println(errors.ErrorStack(err))
		}

		return e, err
	}
	//
	//var events []*event.Event
	//
	//switch r := any(e).(type) {
	//case []*event.Event:
	//	events = r
	//case *event.Event:
	//	if r == nil {
	//		return e, err
	//	}
	//	events = []*event.Event{r}
	//default:
	//	return e, err
	//}

	//output, err := json.Marshal(e)
	//if err != nil {
	//	log.FromContext(ctx.GetContext()).Errorf("error on marshall event for logging. %s", err.Error())
	//} else {
	//	log.FromContext(ctx.GetContext()).WithField("output", output).Info("event sent")
	//}
	//
	//for _, ev := range events {
	//	output, err := json.Marshal(ev)
	//	if err != nil {
	//		log.FromContext(ctx.GetContext()).Errorf("error on marshall event for logging. %s", err.Error())
	//	} else {
	//		log.FromContext(ctx.GetContext()).WithField("output", output).Info("event sent")
	//	}
	//}

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
