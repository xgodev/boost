package nats

import (
	"context"
	"encoding/json"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/nats-io/nats.go"
	"github.com/xgodev/boost/bootstrap/cloudevents"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// SubscriberListener represents a subscriber listener.
type SubscriberListener struct {
	q       *nats.Conn
	handler *cloudevents.HandlerWrapper
	subject string
	queue   string
}

// NewSubscriberListener returns a subscriber listener.
func NewSubscriberListener(q *nats.Conn, handler *cloudevents.HandlerWrapper, subject string,
	queue string) *SubscriberListener {
	return &SubscriberListener{
		q:       q,
		handler: handler,
		subject: subject,
		queue:   queue,
	}
}

// Subscribe subscribes to a particular subject in the listening subscriber's queue.
func (l *SubscriberListener) Subscribe(ctx context.Context) (*nats.Subscription, error) {
	return l.q.QueueSubscribe(l.subject, l.queue, l.h)
}

func (l *SubscriberListener) h(msg *nats.Msg) {

	in := event.New()
	err := json.Unmarshal(msg.Data, &in)
	if err != nil {

		var data interface{}

		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Errorf("could not decode nats record. %s", err.Error())
		} else {
			err := in.SetData("", data)
			if err != nil {
				log.Errorf("could set data from nats record. %s", err.Error())
				return
			}
		}

	}

	logger := log.WithTypeOf(*l).
		WithField("subject", l.subject).
		WithField("queue", l.queue)

	ctx := logger.ToContext(context.Background())

	var inouts []*cloudevents.InOut

	inouts = append(inouts, &cloudevents.InOut{In: &in})

	err = l.handler.Process(ctx, inouts)
	if err != nil {
		logger.Error(errors.ErrorStack(err))
	}

}
