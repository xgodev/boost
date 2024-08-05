package main

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function"
	anats "github.com/xgodev/boost/bootstrap/function/adapter/contrib/nats-io/nats.go/v1"
	lm "github.com/xgodev/boost/bootstrap/function/middleware/logger"
	pm "github.com/xgodev/boost/bootstrap/function/middleware/publisher"
	fnats "github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1"
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/driver/extra/noop"
	"os"
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

func MultiHandler(ctx context.Context, in cloudevents.Event) ([]*cloudevents.Event, error) {
	responseEvent := cloudevents.NewEvent()
	responseEvent.SetID(uuid.New().String())
	responseEvent.SetSource("test-source")
	responseEvent.SetType("test-type")
	responseEvent.SetSubject("test-subject")
	return []*cloudevents.Event{&responseEvent}, nil
}

type H[T any] struct{}

func (h *H[T]) Handle(handler Handler[T]) {
	fmt.Println("Handler executed")
	// Aqui você pode chamar o handler se desejar
	// Exemplo: result, err := handler(context.Background(), event.New())
	// Use result e err conforme necessário
}

/*
func main() {
	h := &H[*cloudevents.Event]{}
	h.Handle(UniqueHandler)
}
*/

func init() {
	os.Setenv("BOOST_FACTORY_ZEROLOG_LEVEL", "TRACE")
	os.Setenv("BOOST_BOOTSTRAP_FUNCTION_DEFAULT", "nats")
}

func main() {

	boost.Start()

	ctx := context.Background()

	conn, err := fnats.NewConn(ctx)
	if err != nil {
		panic(err)
	}

	p := publisher.New(noop.New())
	pmi, err := pm.NewAnyErrorMiddleware[*cloudevents.Event](p)
	if err != nil {
		panic(err)
	}

	lmi, err := lm.NewAnyErrorMiddleware[*cloudevents.Event]()
	if err != nil {
		panic(err)
	}

	fn := function.New[*cloudevents.Event](pmi, lmi)

	err = fn.Run(ctx, UniqueHandler, anats.New[*cloudevents.Event](conn))
	if err != nil {
		panic(err)
	}

}
