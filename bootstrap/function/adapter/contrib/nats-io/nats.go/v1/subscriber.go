package nats

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/nats-io/nats.go"
)

// Subscriber represents a subscriber listener.
type Subscriber[T any] struct {
	q       *nats.Conn
	handler function.Handler[T]
	subject string
	queue   string
}

// NewSubscriber returns a subscriber listener.
func NewSubscriber[T any](q *nats.Conn, handler function.Handler[T], subject string,
	queue string) *Subscriber[T] {
	return &Subscriber[T]{
		q:       q,
		handler: handler,
		subject: subject,
		queue:   queue,
	}
}

// Subscribe subscribes to a particular subject in the listening subscriber's queue.
func (l *Subscriber[T]) Subscribe(ctx context.Context) (*nats.Subscription, error) {
	return l.q.QueueSubscribe(l.subject, l.queue, l.h)
}

func (l *Subscriber[T]) h(msg *nats.Msg) {

	logger := log.WithTypeOf(*l).
		WithField("subject", l.subject).
		WithField("queue", l.queue)

	in := event.New()
	err := json.Unmarshal(msg.Data, &in)
	if err != nil {

		var data interface{}

		if err := json.Unmarshal(msg.Data, &data); err != nil {
			logger.Errorf("could not decode nats record. %s", err.Error())
		} else {
			err := in.SetData("", data)
			if err != nil {
				logger.Errorf("could set data from nats record. %s", err.Error())
				return
			}
		}

	}

	ctx := logger.ToContext(context.Background())

	_, err = l.handler(ctx, in)
	if err != nil {
		logger.Error(errors.ErrorStack(err))
	}

}
