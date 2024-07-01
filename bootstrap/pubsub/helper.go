package pubsub

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/bootstrap/cloudevents"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"

	"cloud.google.com/go/pubsub"
	"github.com/cloudevents/sdk-go/v2/event"
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
		go func(ctx context.Context, m pubsub.Message) {
			h.handle(ctx, m)
		}(ctx, *m)
		m.Ack()
	})
	if err != nil {
		logger.Errorf("pubsub read error: %w", err)
	}
}

func (h *Helper) handle(ctx context.Context, msg pubsub.Message) {

	logger := log.FromContext(ctx).WithTypeOf(*h)

	in := event.New()
	if err := json.Unmarshal(msg.Data, &in); err != nil {
		var data interface{}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			logger.Errorf("could not decode pubsub record. %s", err.Error())
			return
		}

		err := in.SetData("", data)
		if err != nil {
			logger.Errorf("could set data from pubsub record. %s", err.Error())
			return
		}
	}

	var inouts []*cloudevents.InOut

	inouts = append(inouts, &cloudevents.InOut{In: &in})

	if err := h.handler.Process(ctx, inouts); err != nil {
		logger.Error(errors.ErrorStack(err))
	}

}
