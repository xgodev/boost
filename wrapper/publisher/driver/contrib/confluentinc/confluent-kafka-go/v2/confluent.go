package confluent

import (
	"context"
	"encoding/json"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"time"
)

// client represents a Kafka client that implements.
type client struct {
	producer producer
	options  *Options
}

// NewWithConfigPath returns connection with options from config path.
func NewWithConfigPath(ctx context.Context, producer *kafka.Producer, path string) (publisher.Driver, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewWithOptions(ctx, producer, options), nil
}

// NewWithOptions returns connection with options.
func NewWithOptions(ctx context.Context, producer *kafka.Producer, options *Options) publisher.Driver {
	return newDriver(producer, options)
}

func newDriver(producer producer, options *Options) publisher.Driver {
	return &client{producer: producer, options: options}
}

// New returns connection with default options.
func New(ctx context.Context, producer *kafka.Producer) (publisher.Driver, error) {

	options, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewWithOptions(ctx, producer, options), nil
}

// Publish produces every event and waits, bounded by ctx, for the delivery
// report of each one. A conversion error, a Produce error, a failed delivery
// report or the end of ctx before the report arrives becomes the Error of that
// event's output; outputs keep the order of events.
func (p *client) Publish(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Tracef("publishing to kafka %d events", len(events))

	res := make([]publisher.PublishOutput, len(events))
	// Buffered for every event so librdkafka never blocks writing a report
	// after Publish stopped reading because ctx ended.
	deliveryChan := make(chan kafka.Event, len(events))
	pending := make(map[int]struct{}, len(events))

	for i, event := range events {

		res[i].Event = event

		if event.ID() == "" {
			event.SetID(uuid.NewString())
		}

		msg, err := p.convert(ctx, event)
		if err != nil {
			res[i].Error = err
			continue
		}
		msg.Opaque = i

		if err := p.producer.Produce(msg, deliveryChan); err != nil {
			res[i].Error = errors.Wrap(err, errors.Internalf("unable to produce message"))
			continue
		}

		pending[i] = struct{}{}

		p.log(ctx, "message produced, awaiting delivery confirmation")
	}

	p.awaitDeliveries(ctx, deliveryChan, pending, res)

	return res, nil
}

// awaitDeliveries reads delivery reports until every pending output has one or
// ctx ends; then each output still pending gets the ctx error.
func (p *client) awaitDeliveries(ctx context.Context, deliveryChan chan kafka.Event,
	pending map[int]struct{}, res []publisher.PublishOutput) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	for len(pending) > 0 {
		select {
		case e := <-deliveryChan:
			msg, ok := e.(*kafka.Message)
			if !ok {
				logger.Warnf("ignored event on delivery channel: %v", e)
				continue
			}
			i := msg.Opaque.(int) // Publish sets Opaque to the output index
			delete(pending, i)
			if err := msg.TopicPartition.Error; err != nil {
				res[i].Error = errors.Wrap(err, errors.Internalf("unable to deliver message"))
				continue
			}
			p.log(ctx, "message delivered to %s [%d] at offset %v",
				*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
		case <-ctx.Done():
			for i := range pending {
				res[i].Error = errors.Wrap(ctx.Err(), errors.Internalf("delivery report not received"))
			}
			return
		}
	}
}

func (p *client) log(ctx context.Context, format string, args ...interface{}) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	switch p.options.Log.Level {
	case "INFO":
		logger.Infof(format, args...)
	case "TRACE":
		logger.Tracef(format, args...)
	default:
		logger.Debugf(format, args...)
	}
}

func (p *client) convert(ctx context.Context, out *v2.Event) (*kafka.Message, error) {
	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Tracef("converting event to kafka message %s", out.ID())

	var data map[string]interface{}
	if err := out.DataAs(&data); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error on marshal. %s", err.Error()))
	}

	var rawMessage []byte
	rawMessage, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error on marshal. %s", err.Error()))
	}

	pk, err := p.partitionKey(ctx, out)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("unable to gets partition key"))
	}

	headers := p.headers(ctx, out)

	topic := out.Subject()

	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          rawMessage,
		Key:            []byte(pk),
		Headers:        headers,
		Timestamp:      time.Now(),
	}, nil
}

func (p *client) headers(ctx context.Context, out *v2.Event) []kafka.Header {
	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Tracef("setting headers for event %s", out.ID())

	headers := []kafka.Header{
		{Key: "content-type", Value: []byte(out.DataContentType())},
		{Key: "ce_specversion", Value: []byte(out.SpecVersion())},
		{Key: "ce_id", Value: []byte(out.ID())},
		{Key: "ce_source", Value: []byte(out.Source())},
		{Key: "ce_type", Value: []byte(out.Type())},
		{Key: "ce_time", Value: []byte(out.Time().String())},
		{Key: "ce_path", Value: []byte("/")},
		{Key: "ce_subject", Value: []byte(out.Subject())},
	}
	return headers
}

func (p *client) partitionKey(ctx context.Context, out *v2.Event) (string, error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Tracef("getting partition key for event %s", out.ID())

	var pk string
	exts := out.Extensions()

	if key, ok := exts["key"]; ok {
		pk = key.(string)
	} else {
		pk = out.ID()
	}

	return pk, nil
}
