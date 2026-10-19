// Package kafka produces the messages on a topic, through the producer of
// boost's factory.
package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/producer/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2/examples/producer/internal/pkg/domain/publisher/message"
	"github.com/xgodev/boost/model/errors"
)

type publisher struct {
	producer *kafka.Producer
	topic    string
}

// New provides the publisher on the producer of the factory: the generator
// does not build kafka producers, so the constructor calls the factory.
// @Provide(index=0)
func New(ctx context.Context, o *Options) (message.Publisher, error) {
	producer, err := confluent.NewProducer(ctx)
	if err != nil {
		return nil, err
	}
	return &publisher{producer: producer, topic: o.Topic}, nil
}

// Publish produces the message and waits for the broker to take it.
func (p *publisher) Publish(ctx context.Context, msg *entity.Message) error {
	delivered := make(chan kafka.Event, 1)
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg.Text),
	}, delivered)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error producing on topic %s", p.topic))
	}
	select {
	case e := <-delivered:
		if m, ok := e.(*kafka.Message); ok && m.TopicPartition.Error != nil {
			return errors.Wrap(m.TopicPartition.Error, errors.Internalf("error delivering on topic %s", p.topic))
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
