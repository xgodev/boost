package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// ReaderPlugin customizes the reader config before the reader is built (e.g.
// a Dialer with SASL or TLS); kafka.Reader cannot change after it starts.
type ReaderPlugin = factory.Plugin[*kafka.ReaderConfig]

// NewReader returns a new reader with default options.
func NewReader(ctx context.Context, plugins ...ReaderPlugin) (*kafka.Reader, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewReaderWithOptions(ctx, o, plugins...)
}

// NewReaderWithConfigPath returns a new reader with options from config path.
func NewReaderWithConfigPath(ctx context.Context, path string, plugins ...ReaderPlugin) (*kafka.Reader, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewReaderWithOptions(ctx, o, plugins...)
}

// NewReaderWithOptions returns a new reader with options. An invalid config
// is an error here; kafka.NewReader would panic on it.
func NewReaderWithOptions(ctx context.Context, o *Options, plugins ...ReaderPlugin) (*kafka.Reader, error) {

	offset, err := startOffset(o.Reader.StartOffset)
	if err != nil {
		return nil, err
	}

	logger, errorLogger := loggers(ctx)

	c := &kafka.ReaderConfig{
		Brokers:          o.Brokers,
		GroupID:          o.Reader.GroupID,
		Topic:            o.Reader.Topic,
		Partition:        o.Reader.Partition,
		QueueCapacity:    o.Reader.QueueCapacity,
		MinBytes:         o.Reader.MinBytes,
		MaxBytes:         o.Reader.MaxBytes,
		MaxWait:          o.Reader.MaxWait,
		ReadBatchTimeout: o.Reader.ReadBatchTimeout,
		CommitInterval:   o.Reader.CommitInterval,
		StartOffset:      offset,
		Logger:           logger,
		ErrorLogger:      errorLogger,
	}

	if err := factory.ApplyPlugins(ctx, c, plugins...); err != nil {
		return nil, err
	}

	if err := c.Validate(); err != nil {
		return nil, errors.NotValidf("kafka-go reader config: %s", err)
	}

	log.FromContext(ctx).Infof("kafka-go reader created for topic %s on brokers %v", c.Topic, c.Brokers)

	return kafka.NewReader(*c), nil
}
