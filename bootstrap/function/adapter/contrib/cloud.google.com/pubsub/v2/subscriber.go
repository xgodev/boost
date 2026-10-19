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

// Subscribe receives the messages of the subscription until ctx is done, and
// returns once the handlers finished the ones they hold. It returns an error
// naming the subscription when it cannot be received.
func (l *Subscriber[T]) Subscribe(ctx context.Context) error {
	subscription := l.client.Subscriber(l.subscription)
	subscription.ReceiveSettings.MaxOutstandingMessages = int(l.options.Concurrency)

	err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		err := l.processMessage(ctx, msg)

		if err != nil {
			log.Errorf("processing failed: %v", err)
			msg.Nack()

			return
		}

		msg.Ack()
	})
	if err != nil {
		return errors.Wrap(err, errors.Internalf("pubsub: could not receive on %s", l.subscription))
	}

	return nil
}

// processMessage processes each message, retries if needed, and applies backoff
func (l *Subscriber[T]) processMessage(ctx context.Context, msg *pubsub.Message) error {
	logger := log.FromContext(ctx)
	logger.Debugf("processing message %s", msg.ID)

	retryCount := 0

	in, err := generateCloudEvent(msg, l.subscription)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("could not generate CloudEvent: %s", err.Error()))
	}

	for {
		// The stop of the app cancels ctx: the handler finishes the message it
		// holds instead of failing on it.
		if _, err := l.handler(context.WithoutCancel(ctx), in); err != nil {
			retryCount++

			logger.WithField("payload", msg.Data).Warnf("msgID=%s handler failed (attempt %d/%d): %v", msg.ID, retryCount, l.options.RetryLimit, err)
			if retryCount >= l.options.RetryLimit {
				return errors.Wrap(err, errors.Internalf("max retry limit reached"))
			}

			// The app is stopping: give the message back instead of retrying
			if ctx.Err() != nil {
				return errors.Wrap(err, errors.Internalf("stopped before a retry"))
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
