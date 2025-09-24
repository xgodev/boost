package main

import (
	"context"
	"math/rand/v2"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/google/uuid"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	ce "github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http"
	lm "github.com/xgodev/boost/bootstrap/function/middleware/logger"
	pm "github.com/xgodev/boost/bootstrap/function/middleware/publisher"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/driver/extra/noop"
)

func TestLog(ctx context.Context) {

	log.FromContext(ctx).Contextual("testWithoutCtx", "valueN").Info("Test Contextual Log")
	log.FromContext(ctx).WithField("test", "value").Info("Test Message Without Contextual")
	log.FromContext(ctx).Info("test message")
}

func TestLogConcurrency(ctx context.Context, n int) {

	log.FromContext(ctx).Contextual("testWithoutCtx", n).Info("Test Contextual Log")
	time.Sleep(3 * time.Second)
	log.FromContext(ctx).WithField("test", n).Infof("Test Message Without Contextual %d", n)
	log.FromContext(ctx).Info("test message")
}

func Handle(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {
	log.FromContext(ctx).Info("Clean log")

	n := rand.IntN(1000)

	log.FromContext(ctx).Trace("calling external function to include contextual data")
	TestLogConcurrency(ctx, n)
	log.FromContext(ctx).Tracef("end calling for external function %d", n)

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

	p := publisher.New(noop.New())
	pmi, err := pm.NewAnyErrorMiddleware[*cloudevents.Event](p)
	if err != nil {
		panic(err)
	}

	lmi, err := lm.NewAnyErrorMiddleware[*cloudevents.Event]()
	if err != nil {
		panic(err)
	}

	fn, err := function.New[*cloudevents.Event](pmi, lmi)
	if err != nil {
		panic(err)
	}

	err = fn.Run(ctx, Handle, ce.New[*cloudevents.Event](
		[]client.Option{
			cloudevents.WithUUIDs(),
			cloudevents.WithTimeNow(),
		},
	))
	if err != nil {
		panic(err)
	}

}
