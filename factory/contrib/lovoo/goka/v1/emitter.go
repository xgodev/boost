package goka

import (
	"context"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"github.com/xgodev/boost/wrapper/log"
)

// NewEmitterWithConfigPath returns connection with options from config path.
func NewEmitterWithConfigPath(ctx context.Context, path string) (*goka.Emitter, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewEmitterWithOptions(ctx, options)
}

// NewEmitterWithOptions returns connection with options.
func NewEmitterWithOptions(ctx context.Context, o *Options) (*goka.Emitter, error) {

	logger := log.FromContext(ctx)

	topic := goka.Stream(o.Topic)

	emitter, err := goka.NewEmitter(o.Brokers, topic, new(codec.Bytes), goka.WithEmitterLogger(&Logger{}))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}

	logger.Infof("Created kafka connection to %v", o.Brokers)

	return emitter, err
}

// NewEmitter returns connection with default options.
func NewEmitter(ctx context.Context) (*goka.Emitter, error) {

	logger := log.FromContext(ctx)

	o, err := NewOptions()
	if err != nil {
		logger.Fatalf(err.Error())
	}

	return NewEmitterWithOptions(ctx, o)
}
