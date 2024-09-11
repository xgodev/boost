package confluent

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// NewConsumerWithConfigPath returns connection with options from config path.
func NewConsumerWithConfigPath(ctx context.Context, path string) (*kafka.Consumer, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewConsumerWithOptions(ctx, options)
}

// NewConsumerWithOptions returns connection with options.
func NewConsumerWithOptions(ctx context.Context, o *Options) (*kafka.Consumer, error) {

	p, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  o.Brokers,
		"group.id":           o.Consumer.GroupId,
		"auto.offset.reset":  o.Consumer.AutoOffsetReset,
		"enable.auto.commit": o.Consumer.EnableAutoCommit,
		"security.protocol":  o.Consumer.Protocol,
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

// NewConsumer returns connection with default options.
func NewConsumer(ctx context.Context) (*kafka.Consumer, error) {

	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewConsumerWithOptions(ctx, o)
}
