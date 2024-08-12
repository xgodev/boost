package confluent

import (
	"context"
	"encoding/json"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"golang.org/x/sync/errgroup"
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

	if options != nil && options.Log.Enabled {
		logger := NewLogger(producer, options.Log.Level)
		logger.Start()
	}

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
func (p *client) Publish(ctx context.Context, outs []*v2.Event) (err error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Tracef("publishing to kafka %d events", len(outs))

	deliveryChan := make(chan kafka.Event, len(outs))
	var g errgroup.Group

	for _, out := range outs {

		out := out
		g.Go(func() error {
			logger := logger.WithField("subject", out.Subject()).WithField("id", out.ID())

			msg, err := p.convert(ctx, out)
			if err != nil {
				return err
			}

			if err := p.producer.Produce(msg, deliveryChan); err != nil {
				return errors.Wrap(err, errors.Internalf("unable to publish to kafka"))
			}

			logger.Debugf("message produced, awaiting delivery confirmation")
			return nil
		})

	}

	go func() {
		err = g.Wait()
		close(deliveryChan)
	}()

	for e := range deliveryChan {
		m := e.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			return errors.Wrap(m.TopicPartition.Error, errors.Internalf("delivery failed"))
		}

		logger.Infof("message delivered to %s [%d] at offset %v", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	return err
}

func (p *client) convert(ctx context.Context, out *v2.Event) (*kafka.Message, error) {
	var data map[string]interface{}
	if err := out.DataAs(&data); err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error on marshal. %s", err.Error()))
	}

	var rawMessage []byte
	rawMessage, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error on marshal. %s", err.Error()))
	}

	pk, err := p.partitionKey(out)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("unable to gets partition key"))
	}

	headers := p.headers(out)

	topic := out.Subject()

	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          rawMessage,
		Key:            []byte(pk),
		Headers:        headers,
		Timestamp:      time.Now(),
	}, nil
}

func (p *client) headers(out *v2.Event) []kafka.Header {
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

func (p *client) partitionKey(out *v2.Event) (string, error) {

	var pk string
	exts := out.Extensions()

	if key, ok := exts["key"]; ok {
		pk = key.(string)
	} else {
		pk = out.ID()
	}

	return pk, nil
}
