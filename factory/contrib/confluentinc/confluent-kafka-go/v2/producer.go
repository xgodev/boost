package confluent

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
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
			"acks":               o.Producer.Acks,
			"request.timeout.ms": o.Producer.Timeout.Request,
			"message.timeout.ms": o.Producer.Timeout.Message,
		},
	)
	if err != nil {
		return nil, err
	}

	if o.Log.Enabled {
		go func() {
			for e := range p.Events() {

				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						err = errors.Wrap(ev.TopicPartition.Error, errors.Internalf("delivery failed"))
					}
					logger(ctx, o.Log.Level, "message delivered to %s [%d] at offset %v", *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				case kafka.Error:
					log.Errorf("Error: %v\n", ev)
				default:
					log.Warnf("Ignored event: %s\n", ev)
				}
			}
		}()
	}

	return p, nil
}

func logger(ctx context.Context, lvl string, format string, args ...interface{}) {

	l := log.FromContext(ctx)

	switch lvl {
	case "INFO":
		l.Infof(format, args...)
	case "TRACE":
		l.Tracef(format, args...)
	default:
		l.Debugf(format, args...)
	}
}

// NewProducer returns connection with default options.
func NewProducer(ctx context.Context) (*kafka.Producer, error) {

	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewProducerWithOptions(ctx, o)
}
