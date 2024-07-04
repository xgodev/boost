package http

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"
)

func Run(fn function.Handler, opts ...client.Option) error {

	ctx := context.Background()

	logger := log.FromContext(ctx)

	p, err := cloudevents.NewHTTP(cloudevents.WithPort(Port()), cloudevents.WithPath(Path()))
	if err != nil {
		logger.Errorf("failed to create protocol: %s", err.Error())
	}
	c, err := cloudevents.NewClient(p, opts...)
	if err != nil {
		logger.Errorf("failed to create client: %s", err.Error())
		return err
	}

	logger.Infof("listening on :%d%s\n", 8080, "/")

	return c.StartReceiver(ctx, fn)
}
