package pubsub

import (
	"context"
	"math"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
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

	logger.Tracef("pubsub - Subscribing to %s", l.subscription)

	subscription := l.client.Subscriber(l.subscription)
	subscription.ReceiveSettings = pubsub.ReceiveSettings{
		MaxOutstandingMessages: int(l.options.Concurrency),
	}

	err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		err := l.processMessage(ctx, msg)

		if err != nil {
			log.Errorf("processing failed: %v", err)
			msg.Nack()
		}

		msg.Ack()
	})

	if err != nil {
		logger.Fatalf("Failed to start subscription %s: %v", l.subscription, err)
	}

	return nil
}

// processMessage processes each message, retries if needed, and applies backoff
func (l *Subscriber[T]) processMessage(ctx context.Context, msg *pubsub.Message) error {
	logger := log.FromContext(ctx).WithTypeOf(*l)

	retryCount := 0

	in, err := generateCloudEvent(msg, l.subscription)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("could not generate CloudEvent: %s", err.Error()))
	}

	for {
		if _, err := l.handler(ctx, in); err != nil {
			retryCount++

			logger.Warnf("msgID=%s handler failed (attempt %d/%d): %v\nPayload: %s", msg.ID, retryCount, l.options.RetryLimit, err, string(msg.Data))
			if l.options.RetryLimit != -1 && retryCount >= l.options.RetryLimit {
				return errors.Wrap(err, errors.Internalf("max retry limit reached"))
			}

			// Apply backoff if enabled
			if l.options.Backoff {
				l.applyBackoff(retryCount)
			}

			continue
		}

		break
	}

	return nil
}

// applyBackoff applies an exponential backoff strategy
func (l *Subscriber[T]) applyBackoff(retryCount int) {
	backoffTime := time.Duration(math.Pow(2, float64(retryCount))) * l.options.BackoffBase

	// Cap the backoff time
	if backoffTime > l.options.MaxBackoff {
		backoffTime = l.options.MaxBackoff
	}

	time.Sleep(backoffTime)
}
