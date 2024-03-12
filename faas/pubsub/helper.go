package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/xgodev/boost/faas/cloudevents"
	"github.com/xgodev/boost/log"
)

// Helper assists in creating event handlers.
type Helper struct {
	client  *pubsub.Client
	handler *cloudevents.HandlerWrapper
	options *Options
}

// NewHelper returns a new Helper with options.
func NewHelper(ctx context.Context, client *pubsub.Client, options *Options,
	handler *cloudevents.HandlerWrapper) *Helper {

	return &Helper{
		handler: handler,
		options: options,
		client:  client,
	}
}

// NewDefaultHelper returns a new Helper with default options.
func NewDefaultHelper(ctx context.Context, client *pubsub.Client, handler *cloudevents.HandlerWrapper) *Helper {

	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelper(ctx, client, opt, handler)
}

func (h *Helper) Start() {
	topic := h.options.Topic

	go h.run(context.Background(), topic)

	c := make(chan struct{})
	<-c
}

func (h *Helper) run(ctx context.Context, subscriptionName string) {
	logger := log.FromContext(ctx)
	sub := h.client.Subscription(subscriptionName)
	err := sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message: %s", m.Data)
		m.Ack()
	})
	if err != nil {
		logger.Errorf("pubsub read error: %w", err)
	}
}
