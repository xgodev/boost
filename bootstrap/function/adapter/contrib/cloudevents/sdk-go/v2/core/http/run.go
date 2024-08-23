package http

import (
	"context"
	"encoding/json"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"
	"net/http"
)

// Plugin defines a function to process plugin.
type Plugin func(context.Context, []cehttp.Option) []cehttp.Option

func Run[T any](fn function.Handler[T], opts []client.Option, plugins ...Plugin) (err error) {

	ctx := context.Background()

	logger := log.FromContext(ctx)

	httpOpts := []cehttp.Option{ce.WithPort(Port()), ce.WithPath(Path())}

	for _, plugin := range plugins {
		httpOpts = plugin(ctx, httpOpts)
	}

	p, err := ce.NewHTTP(httpOpts...)
	if err != nil {
		logger.Errorf("failed to create protocol: %s", err.Error())
	}

	c, err := ce.NewClient(p, opts...)
	if err != nil {
		logger.Errorf("failed to create client: %s", err.Error())
		return err
	}

	logger.Infof("listening on :%d%s\n", 8080, "/")

	return c.StartReceiver(ctx, Wrapper(fn))
}

func Wrapper[T any](fn function.Handler[T]) func(context.Context, ce.Event) ce.Result {
	return func(ctx context.Context, event ce.Event) ce.Result {
		e, err := fn(ctx, event)
		if err != nil {
			return ce.NewHTTPResult(http.StatusInternalServerError, err.Error())
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
