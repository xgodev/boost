package log

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new logger plugin for nats connection.
func Register(ctx context.Context) (func(context.Context, *nats.Conn) error, func(context.Context, []nats.Option) ([]nats.Option, error)) {
	return nil, handlers
}

// Log represents logger plugin for nats connection.
type Log struct {
	options *Options
}

// NewLogWithOptions returns a new logger plugin with options.
func NewLogWithOptions(options *Options) *Log {
	return &Log{options: options}
}

// NewLogWithConfigPath returns a new logger plugin with options from config path.
func NewLogWithConfigPath(path string) (*Log, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewLogWithOptions(o), nil
}

// NewLog returns a new logger plugin with default options.
func NewLog() *Log {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewLogWithOptions(o)
}

// Register registers this logger plugin for nats connection, unless it is disabled.
func (i *Log) Register(ctx context.Context) (func(context.Context, *nats.Conn) error, func(context.Context, []nats.Option) ([]nats.Option, error)) {
	if !i.options.Enabled {
		return nil, nil
	}
	return nil, handlers
}

func handlers(ctx context.Context, options []nats.Option) ([]nats.Option, error) {
	options = append(options,
		nats.DisconnectErrHandler(disconnectedErrHandler),
		nats.ReconnectHandler(reconnectedHandler),
		nats.ClosedHandler(closedHandler))
	return options, nil
}

func disconnectedErrHandler(nc *nats.Conn, err error) {
	log.Warnf("Disconnected from nats server! will attempt reconnects")
}

func reconnectedHandler(nc *nats.Conn) {
	log.Warnf("Reconnected [%s]", nc.ConnectedUrl())
}

func closedHandler(nc *nats.Conn) {
	log.Errorf("Exiting: %v", nc.LastError())
}
