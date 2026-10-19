package kafka

import (
	"context"
	stderrors "errors"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
)

// writer is the part of *kafka.Writer the driver uses.
type writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type client struct {
	writer writer
}

// New returns a publisher.Driver that writes with w. The writer must not set
// Topic: each event's subject is the topic of its message.
func New(w *kafka.Writer) publisher.Driver {
	return newDriver(w)
}

func newDriver(w writer) publisher.Driver {
	return &client{writer: w}
}

// Publish writes the events in one WriteMessages call, in CloudEvents binary
// mode: the data is the value, the attributes and extensions are ce_*
// headers. The key is the "key" extension, else "partitionkey", else the
// event ID. An event without subject, or one kafka-go reports as failed,
// gets its own Error; outputs keep the order of events.
func (p *client) Publish(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Tracef("publishing to kafka %d events", len(events))

	res := make([]publisher.PublishOutput, len(events))
	msgs := make([]kafka.Message, 0, len(events))
	// sent[j] is the index in events of msgs[j].
	sent := make([]int, 0, len(events))

	for i, event := range events {
		res[i].Event = event

		if event.ID() == "" {
			event.SetID(uuid.NewString())
		}

		msg, err := message(event)
		if err != nil {
			res[i].Error = err
			continue
		}

		msgs = append(msgs, msg)
		sent = append(sent, i)
	}

	if len(msgs) == 0 {
		return res, nil
	}

	err := p.writer.WriteMessages(ctx, msgs...)

	var writeErrors kafka.WriteErrors
	if stderrors.As(err, &writeErrors) && len(writeErrors) == len(msgs) {
		for j, e := range writeErrors {
			if e != nil {
				res[sent[j]].Error = errors.Wrap(e, errors.Internalf("unable to write message to kafka: %s", e))
			}
		}
	} else if err != nil {
		for _, i := range sent {
			res[i].Error = errors.Wrap(err, errors.Internalf("unable to write messages to kafka: %s", err))
		}
	}

	logger.Tracef("published to kafka %d messages", len(msgs))

	return res, nil
}

func message(event *v2.Event) (kafka.Message, error) {

	if event.Subject() == "" {
		return kafka.Message{}, errors.NotValidf("event %s has no subject to use as kafka topic", event.ID())
	}

	h, err := headers(event)
	if err != nil {
		return kafka.Message{}, err
	}

	return kafka.Message{
		Topic:   event.Subject(),
		Key:     []byte(key(event)),
		Value:   event.Data(),
		Headers: h,
	}, nil
}

func key(event *v2.Event) string {
	exts := event.Extensions()
	for _, name := range []string{"key", "partitionkey"} {
		if v, ok := exts[name]; ok {
			if s, err := types.Format(v); err == nil && s != "" {
				return s
			}
		}
	}
	return event.ID()
}

func headers(event *v2.Event) ([]kafka.Header, error) {

	h := []kafka.Header{
		{Key: "ce_specversion", Value: []byte(event.SpecVersion())},
		{Key: "ce_id", Value: []byte(event.ID())},
		{Key: "ce_source", Value: []byte(event.Source())},
		{Key: "ce_type", Value: []byte(event.Type())},
		{Key: "ce_subject", Value: []byte(event.Subject())},
	}

	add := func(key, value string) {
		if value != "" {
			h = append(h, kafka.Header{Key: key, Value: []byte(value)})
		}
	}

	add("content-type", event.DataContentType())
	add("ce_dataschema", event.DataSchema())
	if !event.Time().IsZero() {
		add("ce_time", types.FormatTime(event.Time()))
	}

	for name, v := range event.Extensions() {
		s, err := types.Format(v)
		if err != nil {
			return nil, errors.NotValidf("event %s extension %s: %s", event.ID(), name, err)
		}
		add("ce_"+name, s)
	}

	return h, nil
}
