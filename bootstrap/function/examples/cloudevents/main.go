package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	ce "github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/bootstrap/function/middleware/publisher"
	"github.com/xgodev/boost/bootstrap/function/middleware/publisher/driver/extra/noop"
	"github.com/xgodev/boost/middleware/plugins/local/wrapper/log"
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

	fn := function.New(
		publisher.New(noop.New()),
		log.NewAnyErrorMiddleware[*cloudevents.Event](ctx),
	)

	err := fn.Run(ctx, Handle, ce.New(
		cloudevents.WithUUIDs(),
		cloudevents.WithTimeNow(),
	))
	if err != nil {
		panic(err)
	}

}
