package pubsub

import (
	"context"

	pb "cloud.google.com/go/pubsub/v2"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/log/contrib/rs/zerolog/v1"
)

// Subscriber contains the Pub/Sub client, handler function, and options
type Subscriber[T any] struct {
	client       *pb.Client
	handler      function.Handler[T]
	subscription string
	options      *Options
}

// NewSubscriber returns a subscriber listener.
func NewSubscriber[T any](client *pb.Client, handler function.Handler[T], subscription string, options *Options) *Subscriber[T] {
	return &Subscriber[T]{
		client:       client,
		handler:      handler,
		subscription: subscription,
		options:      options,
	}
}

// Subscribe subscribes and consumes messages from multiple Pub/Sub topics concurrently
func (l *Subscriber[T]) Subscribe(ctx context.Context) error {
	log.Ctx(ctx, *l).Tracef("pubsub - Subscribing to %s", l.subscription)

	subscription := l.client.Subscriber(l.subscription)
	subscription.ReceiveSettings.MaxOutstandingMessages = int(l.options.Concurrency)

	err := subscription.Receive(ctx, func(ctx context.Context, msg *pb.Message) {
		err := l.processMessage(ctx, msg)

		if err != nil {
			msg.Nack()
		}

		msg.Ack()
	})

	if err != nil {
		log.Ctx(ctx, *l).Fatalf("Failed to start subscription %s: %v", l.subscription, err)
	}

	return nil
}

// processMessage processes each message, retries if needed, and applies backoff
func (l *Subscriber[T]) processMessage(ctx context.Context, msg *pb.Message) error {
	ctx = zerolog.NewLogger().ToContext(ctx)

	in, err := generateCloudEvent(msg, l.subscription)
	if err != nil {
		log.Ctx(ctx, *l).Errorf("could not generate CloudEvent: %s", err)
		return errors.Wrap(err, errors.Internalf("could not generate CloudEvent: %s", err.Error()))
	}

	if _, err := l.handler(ctx, in); err != nil {
		a := 1
		if msg.DeliveryAttempt != nil {
			a = *msg.DeliveryAttempt
		}
		log.Ctx(ctx, *l).Warnf("msgID=%s handler failed (attempt %d): %v | Payload: %s", msg.ID, a, err, string(msg.Data))
		return err
	}

	return nil
}
