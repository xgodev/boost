package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

// Subscriber contains the Pub/Sub client and the handler function
type Subscriber[T any] struct {
	client  *pubsub.Client
	handler function.Handler[T]
	topic   string
}

// NewSubscriber returns a subscriber listener.
func NewSubscriber[T any](client *pubsub.Client, handler function.Handler[T], topic string) *Subscriber[T] {
	return &Subscriber[T]{
		topic:   topic,
		handler: handler,
		client:  client,
	}
}

// Subscribe subscribes and consumes messages from multiple Pub/Sub topics concurrently
func (l *Subscriber[T]) Subscribe(ctx context.Context) error {

	logger := log.FromContext(ctx).WithTypeOf(*l)

	// Subscription to the topic
	subscription := l.client.Subscription(fmt.Sprintf("%s-sub", l.topic))

	// Starts the subscription (blocking call in a goroutine)
	err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		logger.Printf("Received message from %s: %s", l.topic, string(msg.Data))

		in := event.New()

		ce := false
		contentType := "application/json"

		// Checks attributes and transforms into a CloudEvent if applicable
		for key, value := range msg.Attributes {
			switch key {
			case "content-type":
				in.SetDataContentType(value)
				contentType = value
			case "ce_specversion":
				in.SetSpecVersion(value)
				ce = true
			case "ce_id":
				in.SetID(value)
				ce = true
			case "ce_source":
				in.SetSource(value)
				ce = true
			case "ce_type":
				in.SetType(value)
				ce = true
			case "ce_time":
				ce = true
				if t, err := time.Parse(time.RFC3339, value); err == nil {
					in.SetTime(t)
				}
			case "ce_subject":
				ce = true
				in.SetSubject(value)
			default:
				in.SetExtension(key, value)
			}
		}

		// If it's not a CloudEvent, create one manually
		if !ce {
			in.SetID(uuid.NewString())
			in.SetSource(fmt.Sprintf("pubsub://%s", l.topic))
			in.SetType("pubsub.message")
			in.SetTime(time.Now())
		}

		// Set the message body as CloudEvent data
		if err := in.SetData(contentType, msg.Data); err != nil {
			logger.Printf("could not set data from pubsub message: %s", err.Error())
		}

		// Processes the event via handler
		if _, err := l.handler(ctx, in); err != nil {
			logger.Printf("Error processing message: %v", err)
			msg.Nack() // Nack the message if there is an error
			return
		}

		// Acknowledge the message after successful processing
		msg.Ack()
	})

	if err != nil {
		logger.Fatalf("Failed to start subscription for topic %s: %v", l.topic, err)
	}

	return nil
}
