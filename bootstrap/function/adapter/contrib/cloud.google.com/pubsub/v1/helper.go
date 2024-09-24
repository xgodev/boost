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
	logger := log.WithTypeOf(*h)
	var wg sync.WaitGroup

	// Subscribe to each topic in a goroutine
	for _, topic := range h.options.Topics {
		wg.Add(1)

		go func(topic string) {
			defer wg.Done()

			subscriber := NewSubscriber[T](h.client, h.handler, topic, h.options)

			// Subscribe to the topic
			if err := subscriber.Subscribe(context.Background()); err != nil {
				logger.Errorf("Failed to subscribe to topic %s: %v", topic, err)
			} else {
				logger.Infof("Successfully subscribed to topic %s", topic)
			}
		}(topic)
	}

	// Wait for all subscriptions to complete
	wg.Wait()
}
