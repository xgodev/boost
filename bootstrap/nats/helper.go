package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/xgodev/boost/bootstrap/cloudevents"
	"github.com/xgodev/boost/wrapper/log"
)

// Helper assists in creating event handlers.
type Helper struct {
	handler  *cloudevents.HandlerWrapper
	queue    string
	subjects []string
	conn     *nats.Conn
}

// NewHelper returns a new Helper with options.
func NewHelper(ctx context.Context, conn *nats.Conn, options *Options,
	handler *cloudevents.HandlerWrapper) *Helper {

	return &Helper{
		handler:  handler,
		queue:    options.Queue,
		subjects: options.Subjects,
		conn:     conn,
	}
}

// NewDefaultHelper returns a new Helper with default options.
func NewDefaultHelper(ctx context.Context, conn *nats.Conn, handler *cloudevents.HandlerWrapper) *Helper {

	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelper(ctx, conn, opt, handler)
}

func (h *Helper) Start() {

	for i := range h.subjects {
		go h.subscribe(context.Background(), h.subjects[i])
	}

	c := make(chan struct{})
	<-c
}

func (h *Helper) subscribe(ctx context.Context, subject string) {

	logger := log.FromContext(ctx).WithTypeOf(*h)

	subscriber := NewSubscriberListener(h.conn, h.handler, subject, h.queue)
	subscribe, err := subscriber.Subscribe(ctx)
	if err != nil {
		logger.Error(err)
	}

	if subscribe.IsValid() {
		logger.Infof("nats: subscribed on %s with queue %s", subject, h.queue)
	} else {
		logger.Errorf("nats: not subscribed on %s with queue %s", subject, h.queue)
	}

}
