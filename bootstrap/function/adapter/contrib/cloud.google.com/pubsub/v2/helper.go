package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub/v2"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"
)

// Helper assists in creating event handlers for Pub/Sub with multiple topics.
type Helper[T any] struct {
	handler function.Handler[T]
	options *Options
	client  *pubsub.Client
}

// NewHelperWithOptions returns a new Helper with custom options.
func NewHelperWithOptions[T any](client *pubsub.Client, handler function.Handler[T], options *Options) *Helper[T] {
	return &Helper[T]{
		handler: handler,
		options: options,
		client:  client,
	}
}

// NewHelper returns a new Helper with default options.
func NewHelper[T any](client *pubsub.Client, handler function.Handler[T]) *Helper[T] {
	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}
	return NewHelperWithOptions(client, handler, opt)
}

// Start subscribes to the topics and processes messages concurrently.
func (h *Helper[T]) Start() {

	// Subscribe to each subscription in a goroutine
	for _, subscription := range h.options.Subscriptions {
		go func(subscription string) {
			// Subscribe to the subscription
			if err := NewSubscriber[T](h.client, h.handler, subscription, h.options).Subscribe(context.Background()); err != nil {
				logger := log.WithTypeOf(*h)
				logger.Errorf("Failed to subscribe to subscription %s: %v", subscription, err)
			}
		}(subscription)
	}
}
