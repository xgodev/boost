package confluent

import (
	"context"
	"encoding/json"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"time"
)

// client represents a Kafka client that implements.
type client struct {
	options *Options
}

// NewWithConfigPath returns connection with options from config path.
func NewWithConfigPath(ctx context.Context, path string) (publisher.Driver, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewWithOptions(ctx, options), nil
}

// NewWithOptions returns connection with options.
func NewWithOptions(ctx context.Context, options *Options) publisher.Driver {
	return &client{options: options}
}

// New returns connection with default options.
func New(ctx context.Context) (publisher.Driver, error) {

	options, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewWithOptions(ctx, options), nil
}

// Publish publishes an event slice.
func (p *client) Publish(ctx context.Context, outs []*v2.Event) (err error) {

	producer, err := confluent.NewProducer(ctx)
	if err != nil {
		return err
	}
	defer producer.Close()

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Tracef("publishing to kafka %d events", len(outs))

	go func() {
		for e := range producer.Events() {

			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					err = errors.Wrap(ev.TopicPartition.Error, errors.Internalf("delivery failed"))
				}
				p.log("message delivered to %s [%d] at offset %v", *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
			case kafka.Error:
				log.Errorf("Error: %v\n", ev)
			default:
				log.Warnf("Ignored event: %s\n", ev)
			}

		}
	}()

	for _, out := range outs {

		if out.ID() == "" {
			out.SetID(uuid.NewString())
		}

		msg, err := p.convert(ctx, out)
		if err != nil {
			return err
		}

		if err := producer.Produce(msg, nil); err != nil {
			return errors.Wrap(err, errors.Internalf("unable to publish to kafka"))
		}

		p.log("message produced, awaiting delivery confirmation")
	}

	for producer.Flush(10000) > 0 {
		p.log("Still waiting to flush outstanding messages")
	}

	return err
}

func (p *client) log(format string, args ...interface{}) {
	switch p.options.Log.Level {
	case "INFO":
		log.Infof(format, args...)
	case "TRACE":
		log.Tracef(format, args...)
	default:
		log.Debugf(format, args...)
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
