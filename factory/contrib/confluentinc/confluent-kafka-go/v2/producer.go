package confluent

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// NewProducerWithConfigPath returns connection with options from config path.
func NewProducerWithConfigPath(ctx context.Context, path string) (*kafka.Producer, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewProducerWithOptions(ctx, options)
}

// NewProducerWithOptions returns connection with options.
func NewProducerWithOptions(ctx context.Context, o *Options) (*kafka.Producer, error) {

	p, err := kafka.NewProducer(
		&kafka.ConfigMap{
			// https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md
			"bootstrap.servers":  o.Brokers,
			"batch.num.messages": o.Producer.Batch.NumMessages,
			"batch.size":         o.Producer.Batch.Size,
			"acks":               o.Acks,
		},
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// NewProducer returns connection with default options.
func NewProducer(ctx context.Context) (*kafka.Producer, error) {

	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewProducerWithOptions(ctx, o)
}
