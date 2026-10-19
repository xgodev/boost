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

// Start runs the helper with no way to stop it, and logs the error Run
// returns.
//
// Deprecated: use Run with the run context, so the app stops on SIGTERM and
// fails when a subscription cannot be received.
func (h *Helper[T]) Start() {
	if err := h.Run(context.Background()); err != nil {
		log.WithTypeOf(*h).Error(err.Error())
	}
}

// Run receives the messages of every subscription concurrently until ctx is
// done. Then it pulls no new message, waits for the handlers to finish the
// ones they hold, and returns. When a subscription cannot be received, for
// instance because it does not exist, Run stops the others and returns that
// error, so the app fails instead of running without it.
func (h *Helper[T]) Run(ctx context.Context) error {
	logger := log.FromContext(ctx).WithTypeOf(*h)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, len(h.options.Subscriptions))
	var wg sync.WaitGroup

	// Subscribe to each subscription in a goroutine
	for _, subscription := range h.options.Subscriptions {
		wg.Add(1)

		go func(subscription string) {
			defer wg.Done()

			logger.Infof("pubsub: receiving on %s", subscription)
			if err := NewSubscriber[T](h.client, h.handler, subscription, h.options).Subscribe(ctx); err != nil {
				errs <- err
				cancel()
			}
		}(subscription)
	}

	// Wait for all subscriptions to stop
	wg.Wait()
	close(errs)
	return <-errs
}
