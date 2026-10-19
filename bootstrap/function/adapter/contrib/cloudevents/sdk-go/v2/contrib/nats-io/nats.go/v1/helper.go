package nats

import (
	"context"
	"sync"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/nats.go"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
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

// Start runs the helper with no way to stop it.
//
// Deprecated: use Run with the run context, so the app stops on SIGTERM.
func (h *Helper[T]) Start() {
	if err := h.Run(context.Background()); err != nil {
		log.WithTypeOf(*h).Error(err)
	}
}

// Run receives the events of every subject until ctx is done. Then each
// receiver drains its subscription, waits for the handlers to finish the
// events they hold, and Run returns. A subject it cannot receive on stops
// the others and fails Run.
func (h *Helper[T]) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, len(h.subjects))
	var wg sync.WaitGroup
	for _, subject := range h.subjects {
		wg.Add(1)
		go func(subject string) {
			defer wg.Done()
			if err := h.subscribe(ctx, subject); err != nil {
				errs <- err
				cancel()
			}
		}(subject)
	}
	wg.Wait()
	close(errs)
	return <-errs
}

func (h *Helper[T]) subscribe(ctx context.Context, subject string) error {

	p, err := cenats.NewConsumerFromConn(h.conn, subject, cenats.WithQueueSubscriber(h.queue))
	if err != nil {
		return errors.Wrap(err, errors.Internalf("cloudevents nats: could not create the consumer of %s", subject))
	}

	defer p.Close(ctx)

	c, err := cloudevents.NewClient(p)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("cloudevents nats: could not create the client of %s", subject))
	}

	// the nats consumer has no responder, so the client refuses a receiver
	// that returns an event: what the function returns is for its
	// middlewares, like publisher, not for the protocol. The stop of the app
	// cancels ctx: the handler finishes the event it holds instead of
	// failing on it.
	receive := func(ctx context.Context, in cloudevents.Event) error {
		_, err := h.handler(context.WithoutCancel(ctx), in)
		return err
	}
	log.FromContext(ctx).WithTypeOf(*h).Infof("cloudevents nats: receiving on %s with queue %s", subject, h.queue)
	if err := c.StartReceiver(ctx, receive); err != nil {
		return errors.Wrap(err, errors.Internalf("cloudevents nats: could not receive on %s with queue %s", subject, h.queue))
	}
	return nil
}
