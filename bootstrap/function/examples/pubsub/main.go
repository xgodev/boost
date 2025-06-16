package main

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	apubsub "github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloud.google.com/pubsub/v1"
	lm "github.com/xgodev/boost/bootstrap/function/middleware/logger"
	pm "github.com/xgodev/boost/bootstrap/function/middleware/publisher"
	fpubsub "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1"
	"github.com/xgodev/boost/wrapper/publisher"
	drvpubsub "github.com/xgodev/boost/wrapper/publisher/driver/contrib/cloud.google.com/pubsub/v1"
)

type Handler[T any] func(context.Context, cloudevents.Event) (T, error)

func UniqueHandler(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {
	responseEvent := cloudevents.NewEvent()
	responseEvent.SetID(uuid.New().String())
	responseEvent.SetSource("test-source")
	responseEvent.SetType("test-type")
	responseEvent.SetSubject("test-subject")
	return &responseEvent, nil
}

type H[T any] struct{}

func (h *H[T]) Handle(handler Handler[T]) {
	fmt.Println("Handler executed")
}

func main() {

	boost.Start()

	ctx := context.Background()

	pb, err := fpubsub.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	dpubsub, err := drvpubsub.New(ctx, pb)
	if err != nil {
		panic(err)
	}

	p := publisher.New(dpubsub)

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

	adpt := apubsub.New[*cloudevents.Event](pb)

	err = fn.Run(ctx, UniqueHandler, adpt)
	if err != nil {
		panic(err)
	}

}
