package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/xgodev/boost"
	fpubsub "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	drvpubsub "github.com/xgodev/boost/wrapper/publisher/driver/contrib/cloud.google.com/pubsub/v1"

	v2 "github.com/cloudevents/sdk-go/v2"
)

func main() {
	boost.Start()

	ctx := context.Background()
	logger := log.FromContext(ctx)

	// Create pubsub v2 client
	client, err := fpubsub.NewClient(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Create publisher driver
	driver, err := drvpubsub.New(ctx, client)
	if err != nil {
		logger.Fatal(err.Error())
	}

	p := publisher.New(driver)

	// Build a CloudEvent
	event := v2.NewEvent()
	event.SetID(uuid.New().String())
	event.SetSource("example/source")
	event.SetType("example.type")
	event.SetSubject("test-topic")
	if err := event.SetData("application/json", map[string]string{"message": "Hello, Pub/Sub v2!"}); err != nil {
		logger.Fatal(err.Error())
	}

	// Publish the event
	if err := p.Publish(ctx, []*v2.Event{&event}); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Infof("event published successfully")
}