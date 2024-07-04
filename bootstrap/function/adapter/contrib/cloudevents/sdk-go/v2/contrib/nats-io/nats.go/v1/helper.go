package nats

import (
	"context"
	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/nats.go"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/log"
)

// Helper assists in creating event handlers.
type Helper struct {
	handler  function.Handler
	queue    string
	subjects []string
	conn     *nats.Conn
}

// NewHelper returns a new Helper with options.
func NewHelper(conn *nats.Conn, options *Options,
	handler function.Handler) *Helper {

	return &Helper{
		handler:  handler,
		queue:    options.Queue,
		subjects: options.Subjects,
		conn:     conn,
	}
}

// NewDefaultHelper returns a new Helper with default options.
func NewDefaultHelper(conn *nats.Conn, handler function.Handler) *Helper {

	opt, err := DefaultOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelper(conn, opt, handler)
}

func (h *Helper) Start() {

	for i := range h.subjects {
		go h.subscribe(context.Background(), h.subjects[i])
	}

	c := make(chan struct{})
	<-c
}

func (h *Helper) subscribe(ctx context.Context, subject string) {

	logger := log.FromContext(ctx)

	p, err := cenats.NewConsumerFromConn(h.conn, subject)
	if err != nil {
		logger.Fatalf("failed to create nats protocol, %s", err.Error())
	}

	defer p.Close(ctx)

	c, err := cloudevents.NewClient(p)
	if err != nil {
		logger.Fatalf("failed to create client, %s", err.Error())
	}

	for {
		if err := c.StartReceiver(ctx, h.handler); err != nil {
			logger.Printf("failed to start nats receiver, %s", err.Error())
		}
	}

}
