package datadog

import (
	"context"

	datadog "github.com/xgodev/boost/factory/contrib/datadog/dd-trace-go/v1"
	"github.com/xgodev/boost/wrapper/log"
	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// Register registers a new datadog plugin for grpc client.
func Register(ctx context.Context) ([]grpc.DialOption, []grpc.CallOption) {
	o, err := NewOptions()
	if err != nil {
		return nil, nil
	}
	h := NewDatadogWithOptions(o)
	return h.Register(ctx)
}

// Datadog represents datadog plugin for grpc client.
type Datadog struct {
	options *Options
}

// NewDatadogWithOptions returns a new datadog plugin with options.
func NewDatadogWithOptions(options *Options) *Datadog {
	return &Datadog{options: options}
}

// NewDatadogWithConfigPath returns a new datadog plugin with options from config path.
func NewDatadogWithConfigPath(path string, traceOptions ...grpctrace.Option) (*Datadog, error) {
	o, err := NewOptionsWithPath(path, traceOptions...)
	if err != nil {
		return nil, err
	}
	return NewDatadogWithOptions(o), nil
}

// NewDatadog returns a new datadog plugin with default options.
func NewDatadog(traceOptions ...grpctrace.Option) *Datadog {
	o, err := NewOptions(traceOptions...)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewDatadogWithOptions(o)
}

// Register registers this datadog plugin for grpc client.
func (i *Datadog) Register(ctx context.Context) ([]grpc.DialOption, []grpc.CallOption) {

	if !i.options.Enabled || !datadog.IsTracerEnabled() {
		return nil, nil
	}

	logger := log.FromContext(ctx)
	logger.Debug("datadog interceptor successfully enabled in grpc client")

	return []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(grpctrace.UnaryClientInterceptor(i.options.traceOptions...)),
		grpc.WithChainStreamInterceptor(grpctrace.StreamClientInterceptor(i.options.traceOptions...)),
	}, nil

}
