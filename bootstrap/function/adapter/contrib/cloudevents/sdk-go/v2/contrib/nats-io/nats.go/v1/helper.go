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
type Helper[T any] struct {
	handler  function.Handler[T]
	queue    string
	subjects []string
	conn     *nats.Conn
}

// NewHelperWithOptions returns a new Helper with options.
func NewHelperWithOptions[T any](conn *nats.Conn,
	handler function.Handler[T], options *Options) *Helper[T] {

	return &Helper[T]{
		handler:  handler,
		queue:    options.Queue,
		subjects: options.Subjects,
		conn:     conn,
	}
}

// NewHelper returns a new Helper with default options.
func NewHelper[T any](conn *nats.Conn, handler function.Handler[T]) *Helper[T] {

	opt, err := NewOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	return NewHelperWithOptions(conn, handler, opt)
}

func (h *Helper[T]) Start() {

	for i := range h.subjects {
		go h.subscribe(context.Background(), h.subjects[i])
	}

	c := make(chan struct{})
	<-c
}

func (h *Helper[T]) subscribe(ctx context.Context, subject string) {

	logger := log.FromContext(ctx)

	p, err := cenats.NewConsumerFromConn(h.conn, subject, cenats.WithQueueSubscriber(h.queue))
	if err != nil {
		logger.Fatalf("failed to create nats protocol, %s", err.Error())
	}

	defer p.Close(ctx)

	c, err := cloudevents.NewClient(p)
	if err != nil {
		logger.Fatalf("failed to create client, %s", err.Error())
		return
	}

	if err := c.StartReceiver(ctx, h.handler); err != nil {
		logger.Printf("failed to start nats receiver, %s", err.Error())
	}

}
