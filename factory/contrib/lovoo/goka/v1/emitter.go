package goka

import (
	"context"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"github.com/xgodev/boost/wrapper/log"
)

// NewEmitterWithConfigPath returns connection with options from config path.
func NewEmitterWithConfigPath(ctx context.Context, path string) (*Emitter, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewEmitterWithOptions(ctx, options), nil
}

// NewEmitterWithOptions returns connection with options.
func NewEmitterWithOptions(ctx context.Context, o *Options) *Emitter {
	return &Emitter{options: o}
}

// NewEmitter returns connection with default options.
func NewEmitter(ctx context.Context) (*Emitter, error) {

	logger := log.FromContext(ctx)

	o, err := NewOptions()
	if err != nil {
		logger.Fatalf(err.Error())
	}

	return NewEmitterWithOptions(ctx, o), nil
}

type Emitter struct {
	options *Options
}

// EmitWithHeaders sends a message with the given headers for the passed key using the emitter's codec.
func (e *Emitter) EmitWithHeaders(ctx context.Context, topic string, key string, msg interface{}, headers goka.Headers) (*goka.Promise, error) {
	ge, err := newEmitter(ctx, e.options, topic)
	if err != nil {
		return nil, err
	}
	return ge.EmitWithHeaders(key, msg, headers)
}

// Emit sends a message for passed key using the emitter's codec.
func (e *Emitter) Emit(ctx context.Context, topic string, key string, msg interface{}) (*goka.Promise, error) {
	return e.EmitWithHeaders(ctx, topic, key, msg, nil)
}

// EmitSyncWithHeaders sends a message with the given headers to passed topic and key.
func (e *Emitter) EmitSyncWithHeaders(ctx context.Context, topic string, key string, msg interface{}, headers goka.Headers) error {
	ge, err := newEmitter(ctx, e.options, topic)
	if err != nil {
		return err
	}
	return ge.EmitSyncWithHeaders(key, msg, headers)
}

// EmitSync sends a message to passed topic and key.
func (e *Emitter) EmitSync(ctx context.Context, topic string, key string, msg interface{}) error {
	return e.EmitSyncWithHeaders(ctx, topic, key, msg, nil)
}

func newEmitter(ctx context.Context, o *Options, topic string) (*goka.Emitter, error) {

	logger := log.FromContext(ctx)

	t := goka.Stream(topic)

	emitter, err := goka.NewEmitter(o.Brokers, t, new(codec.Bytes), goka.WithEmitterLogger(&Logger{}))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}

	logger.Infof("Created kafka connection to %v", o.Brokers)

	return emitter, err
}
