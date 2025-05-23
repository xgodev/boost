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
	producer *kafka.Producer
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

// Publish publishes an event slice.
func (p *client) Publish(ctx context.Context, outs []*v2.Event) (res []publisher.PublishOutput, err error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Tracef("publishing to kafka %d events", len(outs))

	for _, out := range outs {

		if out.ID() == "" {
			out.SetID(uuid.NewString())
		}

		msg, err := p.convert(ctx, out)
		if err != nil {
			res = append(res, publisher.PublishOutput{Event: out, Error: err})
			continue
		}

		if err := p.producer.Produce(msg, nil); err != nil {
			res = append(res, publisher.PublishOutput{Event: out, Error: errors.Wrap(err, errors.Internalf("unable to produce message"))})
			continue
		}

		res = append(res, publisher.PublishOutput{Event: out})

		p.log(ctx, "message produced, awaiting delivery confirmation")
	}

	for p.producer.Flush(10000) > 0 {
		p.log(ctx, "Still waiting to flush outstanding messages")
	}

	return res, err
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
