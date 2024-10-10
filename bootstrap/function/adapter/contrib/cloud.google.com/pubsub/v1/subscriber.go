package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"math"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

// Subscriber contains the Pub/Sub client, handler function, and options
type Subscriber[T any] struct {
	client       *pubsub.Client
	handler      function.Handler[T]
	subscription string
	options      *Options
}

// NewSubscriber returns a subscriber listener.
func NewSubscriber[T any](client *pubsub.Client, handler function.Handler[T], subscription string, options *Options) *Subscriber[T] {
	return &Subscriber[T]{
		client:       client,
		handler:      handler,
		subscription: subscription,
		options:      options,
	}
}

// Subscribe subscribes and consumes messages from multiple Pub/Sub topics concurrently
func (l *Subscriber[T]) Subscribe(ctx context.Context) error {
	logger := log.FromContext(ctx).WithTypeOf(*l)

	subscription := l.client.Subscription(l.subscription)
	subscription.ReceiveSettings = pubsub.ReceiveSettings{
		MaxOutstandingMessages: int(l.options.Concurrency),
	}

	err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Println("Received message")
		err := l.processMessage(ctx, msg)

		if err != nil {
			log.Errorf(err.Error())
			msg.Nack()
		}
	})

	if err != nil {
		logger.Fatalf("Failed to start subscription %s: %v", l.subscription, err)
	}

	return nil
}

// processMessage processes each message, retries if needed, and applies backoff
func (l *Subscriber[T]) processMessage(ctx context.Context, msg *pubsub.Message) error {
	retryCount := 0

	in, err := l.generateCloudEvent(msg)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("could not generate CloudEvent: %s", err.Error()))
	}

	for {
		// Processes the event via handler
		if _, err := l.handler(ctx, in); err != nil {
			retryCount++

			// Check retry limit
			if l.options.RetryLimit != -1 && retryCount >= l.options.RetryLimit {
				return errors.Wrap(err, errors.Internalf("max retry limit reached"))
			}

			// Apply backoff if enabled
			if l.options.Backoff {
				l.applyBackoff(retryCount)
			}

			// Retry processing the message
			continue
		}

		// Acknowledge the message after successful processing
		msg.Ack()
		break
	}

	return nil
}

func (l *Subscriber[T]) generateCloudEvent(msg *pubsub.Message) (event.Event, error) {
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
		in.SetSource(fmt.Sprintf("pubsub://%s", l.subscription))
		in.SetType("pubsub.message")
		in.SetTime(time.Now())
	}

	// Set the message body as CloudEvent data
	if err := in.SetData(contentType, msg.Data); err != nil {
		return event.Event{}, errors.Wrap(err, errors.Internalf("could not set data from pubsub message: %s", err.Error()))
	}
	return in, nil
}

// applyBackoff applies an exponential backoff strategy
func (l *Subscriber[T]) applyBackoff(retryCount int) {
	backoffTime := time.Duration(math.Pow(2, float64(retryCount))) * l.options.BackoffBase

	// Cap the backoff time
	if backoffTime > l.options.MaxBackoff {
		backoffTime = l.options.MaxBackoff
	}
}
