package http

import (
	"context"
	"encoding/json"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/go-playground/validator/v10"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	response "github.com/xgodev/boost/model/restresponse"
	"github.com/xgodev/boost/wrapper/log"
	"net/http"
)

// Plugin defines a function to process plugin.
type Plugin func(context.Context, []cehttp.Option) []cehttp.Option

// Run serves fn with a context that is never cancelled, so it never returns
// on its own.
//
// Deprecated: use RunContext, which stops when its context is done.
func Run[T any](fn function.Handler[T], opts []client.Option, plugins ...Plugin) error {
	return RunContext[T](context.Background(), fn, opts, plugins...)
}

// RunContext serves fn over HTTP on Port() and Path() until ctx is done:
// then it takes no new event, shuts the server down once the events in hand
// are answered, and returns. A handler runs with a context the stop does not
// cancel, so the event it holds finishes.
func RunContext[T any](ctx context.Context, fn function.Handler[T], opts []client.Option, plugins ...Plugin) error {
	logger := log.FromContext(ctx)

	httpOpts := []cehttp.Option{ce.WithPort(Port()), ce.WithPath(Path())}

	for _, plugin := range plugins {
		httpOpts = plugin(ctx, httpOpts)
	}

	p, err := ce.NewHTTP(httpOpts...)
	if err != nil {
		logger.Errorf("failed to create protocol: %s", err.Error())
		return err
	}

	c, err := ce.NewClient(p, opts...)
	if err != nil {
		logger.Errorf("failed to create client: %s", err.Error())
		return err
	}

	logger.Infof("listening on :%d%s", Port(), Path())

	handle := Wrapper(fn)
	return c.StartReceiver(ctx, func(ctx context.Context, event ce.Event) ce.Result {
		return handle(context.WithoutCancel(ctx), event)
	})
}

func Wrapper[T any](fn function.Handler[T]) func(context.Context, ce.Event) ce.Result {
	return func(ctx context.Context, event ce.Event) ce.Result {
		e, err := fn(ctx, event)
		if err != nil {
			if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreAsSuccess != 0 {
				return ce.NewHTTPResult(http.StatusOK, "")
			}
			status := ErrorStatusCode(err)
			return ce.NewHTTPResult(status, err.Error())
		}

		var result []byte

		switch v := any(e).(type) {
		case []*ce.Event:
			if len(v) > 0 {
				result, err = json.Marshal(v)
				if err != nil {
					return ce.NewHTTPResult(http.StatusInternalServerError, err.Error())
				}
				return ce.NewHTTPResult(http.StatusOK, string(result))
			}
			return ce.NewHTTPResult(http.StatusOK, "[]")
		case *ce.Event:
			result, err = json.Marshal(v)
			if err != nil {
				return ce.NewHTTPResult(http.StatusInternalServerError, err.Error())
			}
			return ce.NewHTTPResult(http.StatusOK, string(result))
		default:
			return ce.NewHTTPResult(http.StatusInternalServerError, "unsupported handler type")
		}
	}
}

// ErrorStatusCode translates err to the respective HTTP status code.
func ErrorStatusCode(err error) int {
	if _, ok := err.(validator.ValidationErrors); ok {
		return http.StatusUnprocessableEntity
	}
	return response.HTTPStatusFor(errors.Classify(err))
}
