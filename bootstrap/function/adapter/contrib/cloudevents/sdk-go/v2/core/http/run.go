package http

import (
	"context"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"
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

	return c.StartReceiver(ctx, fn)
}
