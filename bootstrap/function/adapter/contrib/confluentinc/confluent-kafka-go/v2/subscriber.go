package confluent

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
)

// Subscriber represents a subscriber listener.
type Subscriber[T any] struct {
	consumer *kafka.Consumer
	handler  function.Handler[T]
	topics   []string
	timeOut  time.Duration
}

// NewSubscriber returns a subscriber listener.
func NewSubscriber[T any](consumer *kafka.Consumer, handler function.Handler[T], topics []string, timeOut time.Duration) *Subscriber[T] {
	return &Subscriber[T]{
		consumer: consumer,
		handler:  handler,
		topics:   topics,
		timeOut:  timeOut,
	}
}

// Subscribe subscribes to a particular subject in the listening subscriber's queue.
func (l *Subscriber[T]) Subscribe(ctx context.Context) error {

	logger := log.FromContext(ctx)

	if err := l.consumer.SubscribeTopics(l.topics, nil); err != nil {
		return err
	}

	run := true
	for run {
		msg, err := l.consumer.ReadMessage(l.timeOut)
		if err != nil {
			if err.(kafka.Error).IsTimeout() {
				logger.Warnf("Consumer error: %v (%v)", err, msg)
				continue
			}
			return err
		}

		logger.Tracef("Message on %s: %s", msg.TopicPartition, string(msg.Value))

		in := event.New()

		ce := false
		contentType := "application/json"

		if msg.Headers != nil {
			for _, h := range msg.Headers {

				switch h.Key {
				case "content-type":
					in.SetDataContentType(string(h.Value))
					contentType = string(h.Value)
				case "ce_specversion":
					in.SetSpecVersion(string(h.Value))
					ce = true
				case "ce_id":
					in.SetID(string(h.Value))
					ce = true
				case "ce_source":
					in.SetSource(string(h.Value))
				case "ce_type":
					in.SetType(string(h.Value))
				case "ce_time":
					if t, err := time.Parse(time.RFC3339, string(h.Value)); err != nil {
						in.SetTime(t)
					}
				case "ce_subject":
					in.SetSubject(string(h.Value))
				default:
					in.SetExtension(h.Key, string(h.Value))
				}
			}
		}

		if !ce {
			in.SetID(uuid.NewString())
			// TODO: adds another default values
		}

		if err := in.SetData(contentType, msg.Value); err != nil {
			logger.Warnf("could not set data from kafka record. %s", err.Error())
		}

		_, err = l.handler(ctx, in)
		if err != nil {
			logger.Error(errors.ErrorStack(err))
		}

	}

	return nil
}
