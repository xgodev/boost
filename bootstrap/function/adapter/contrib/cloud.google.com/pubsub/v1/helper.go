package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"
	"sync"
)

// Helper assists in creating event handlers for Pub/Sub with multiple topics.
type Helper[T any] struct {
	handler function.Handler[T]
	topics  []string
	client  *pubsub.Client
}

// NewHelperWithOptions returns a new Helper with custom options.
func NewHelperWithOptions[T any](client *pubsub.Client, handler function.Handler[T], topics []string) *Helper[T] {
	return &Helper[T]{
		handler: handler,
		topics:  topics,
		client:  client,
	}
}

// NewHelper returns a new Helper with default options.
func NewHelper[T any](client *pubsub.Client, handler function.Handler[T]) *Helper[T] {
	return &Helper[T]{
		handler: handler,
		client:  client,
	}
}

// Start subscribes to the topics and processes messages concurrently.
func (h *Helper[T]) Start(ctx context.Context) {
	logger := log.FromContext(ctx).WithTypeOf(*h)
	var wg sync.WaitGroup

	// Subscribe to each topic in a goroutine
	for _, topic := range h.topics {
		wg.Add(1)

		go func(topic string) {
			defer wg.Done()

			subscriber := NewSubscriber[T](h.client, h.handler, topic)

			// Subscribe to the topic
			if err := subscriber.Subscribe(ctx); err != nil {
				logger.Errorf("Failed to subscribe to topic %s: %v", topic, err)
			} else {
				logger.Infof("Successfully subscribed to topic %s", topic)
			}
		}(topic)
	}

	// Wait for all subscriptions to complete
	wg.Wait()
}
