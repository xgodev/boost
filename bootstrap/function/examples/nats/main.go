package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	anats "github.com/xgodev/boost/bootstrap/function/adapter/contrib/nats-io/nats.go/v1"
	"github.com/xgodev/boost/bootstrap/function/middleware/publisher"
	"github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/log"
	fnats "github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1"
	"os"
)

func Handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {

	responseEvent := cloudevents.NewEvent()
	responseEvent.SetID(uuid.New().String())
	responseEvent.SetSource("test-source")
	responseEvent.SetType("test-type")
	responseEvent.SetSubject("test-subject")
	return &responseEvent, nil

}

func init() {
	os.Setenv("BOOST_FACTORY_ZEROLOG_LEVEL", "TRACE")
}

func main() {

	boost.Start()

	ctx := context.Background()

	conn, err := fnats.NewConn(ctx)
	if err != nil {
		panic(err)
	}

	fn := function.New(
		publisher.New(nats.New(conn)),
		log.NewAnyErrorMiddleware[*cloudevents.Event](ctx),
	)

	err = fn.Run(ctx, Handle, anats.New(conn))
	if err != nil {
		panic(err)
	}

}
