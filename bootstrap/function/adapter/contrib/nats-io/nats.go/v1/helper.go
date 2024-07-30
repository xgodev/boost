package nats

import (
	"context"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"

	"github.com/nats-io/nats.go"
)

// Helper assists in creating event handlers.
type Helper struct {
	handler  function.Handler
	queue    string
	subjects []string
	conn     *nats.Conn
}

// NewHelperWithOptions returns a new Helper with options.
func NewHelperWithOptions(conn *nats.Conn, handler function.Handler, options *Options) *Helper {

	return &Helper{
		handler:  handler,
		queue:    options.Queue,
		subjects: options.Subjects,
		conn:     conn,
	}
}

// NewHelper returns a new Helper with default options.
func NewHelper(conn *nats.Conn, handler function.Handler) *Helper {

	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelperWithOptions(conn, handler, opt)
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

	subscriber := NewSubscriber(h.conn, h.handler, subject, h.queue)
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
