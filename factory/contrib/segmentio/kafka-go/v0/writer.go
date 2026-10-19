package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// WriterPlugin customizes the writer before its first use (e.g. a Transport
// with SASL or TLS).
type WriterPlugin = factory.Plugin[*kafka.Writer]

// NewWriter returns a new writer with default options.
func NewWriter(ctx context.Context, plugins ...WriterPlugin) (*kafka.Writer, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewWriterWithOptions(ctx, o, plugins...)
}

// NewWriterWithConfigPath returns a new writer with options from config path.
func NewWriterWithConfigPath(ctx context.Context, path string, plugins ...WriterPlugin) (*kafka.Writer, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewWriterWithOptions(ctx, o, plugins...)
}

// NewWriterWithOptions returns a new writer with options. The writer has no
// Topic: each message names its own. It connects on the first write.
func NewWriterWithOptions(ctx context.Context, o *Options, plugins ...WriterPlugin) (*kafka.Writer, error) {

	if len(o.Brokers) == 0 {
		return nil, errors.NotValidf("kafka-go brokers are empty")
	}

	b, err := balancer(o.Writer.Balancer)
	if err != nil {
		return nil, err
	}
	acks, err := requiredAcks(o.Writer.RequiredAcks)
	if err != nil {
		return nil, err
	}
	c, err := compression(o.Writer.Compression)
	if err != nil {
		return nil, err
	}

	logger, errorLogger := loggers(ctx)

	w := &kafka.Writer{
		Addr:                   kafka.TCP(o.Brokers...),
		Balancer:               b,
		RequiredAcks:           acks,
		Compression:            c,
		MaxAttempts:            o.Writer.MaxAttempts,
		BatchSize:              o.Writer.BatchSize,
		BatchBytes:             o.Writer.BatchBytes,
		BatchTimeout:           o.Writer.BatchTimeout,
		ReadTimeout:            o.Writer.ReadTimeout,
		WriteTimeout:           o.Writer.WriteTimeout,
		Async:                  o.Writer.Async,
		AllowAutoTopicCreation: o.Writer.AllowAutoTopicCreation,
		Logger:                 logger,
		ErrorLogger:            errorLogger,
	}

	if err := factory.ApplyPlugins(ctx, w, plugins...); err != nil {
		return nil, err
	}

	log.FromContext(ctx).Infof("kafka-go writer created for brokers %v", o.Brokers)

	return w, nil
}
