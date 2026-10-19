package nats

import (
	"context"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"

	"github.com/nats-io/nats.go"
)

// Helper assists in creating event handlers.
type Helper[T any] struct {
	handler  function.Handler[T]
	queue    string
	subjects []string
	conn     *nats.Conn
}

// NewHelperWithOptions returns a new Helper with options.
func NewHelperWithOptions[T any](conn *nats.Conn, handler function.Handler[T], options *Options) *Helper[T] {

	return &Helper[T]{
		handler:  handler,
		queue:    options.Queue,
		subjects: options.Subjects,
		conn:     conn,
	}
}

// NewHelper returns a new Helper with default options.
func NewHelper[T any](conn *nats.Conn, handler function.Handler[T]) *Helper[T] {

	opt, err := DefaultOptions()
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

// Run subscribes to every subject and handles their messages until ctx is
// done. Then it drains the subscriptions: it takes no new message, lets the
// handlers finish the ones they hold, and returns. A subject it cannot
// subscribe to fails Run, after unsubscribing the ones it already has.
func (h *Helper[T]) Run(ctx context.Context) error {
	logger := log.FromContext(ctx).WithTypeOf(*h)

	var subs []*nats.Subscription
	for _, subject := range h.subjects {
		sub, err := NewSubscriber[T](h.conn, h.handler, subject, h.queue).Subscribe(ctx)
		if err != nil {
			for _, s := range subs {
				_ = s.Unsubscribe()
			}
			return errors.Wrap(err, errors.Internalf("nats: could not subscribe on %s with queue %s", subject, h.queue))
		}
		logger.Infof("nats: subscribed on %s with queue %s", subject, h.queue)
		subs = append(subs, sub)
	}

	<-ctx.Done()

	// drain them all before waiting, so no subject takes new messages while
	// another one finishes
	var closed []<-chan nats.SubStatus
	for _, sub := range subs {
		ch := sub.StatusChanged(nats.SubscriptionClosed)
		if err := sub.Drain(); err != nil {
			logger.Warnf("nats: could not drain %s: %v", sub.Subject, err)
			continue
		}
		closed = append(closed, ch)
	}
	for _, ch := range closed {
		<-ch
	}
	return nil
}
