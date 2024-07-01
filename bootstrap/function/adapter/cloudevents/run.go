package cloudevents

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/xgodev/boost/wrapper/log"
)

func Run(fn interface{}, opts ...client.Option) error {

	ctx := context.Background()

	p, err := cloudevents.NewHTTP(cloudevents.WithPort(Port()), cloudevents.WithPath(Path()))
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}
	c, err := cloudevents.NewClient(p, opts...)
	if err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	log.Printf("listening on :%d%s\n", 8080, "/")

	return c.StartReceiver(ctx, fn)
}
