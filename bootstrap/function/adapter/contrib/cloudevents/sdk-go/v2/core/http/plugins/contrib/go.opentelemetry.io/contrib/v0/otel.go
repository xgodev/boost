package contrib

import (
	"context"
	ce "github.com/cloudevents/sdk-go/v2"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"

	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new opentelemetry plugin for http cloudevents.
func Register(ctx context.Context, opts []cehttp.Option) []cehttp.Option {
	o, err := NewOptions()
	if err != nil {
		return nil
	}
	h := NewOtelWithOptions(o)
	return h.Register(ctx, opts)
}

// Otel represents opentelemetry plugin for http cloudevents.
type Otel struct {
	options *Options
}

// NewOtelWithOptions returns a new opentelemetry plugin with options.
func NewOtelWithOptions(options *Options) *Otel {
	return &Otel{options: options}
}

// NewOtelWithConfigPath returns a new opentelemetry plugin with options from config path.
func NewOtelWithConfigPath(path string) (*Otel, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewOtelWithOptions(o), nil
}

// NewOtel returns a new opentelemetry plugin with default options.
func NewOtel() *Otel {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewOtelWithOptions(o)
}

// Register registers this opentelemetry plugin for http cloudevents.
func (i *Otel) Register(ctx context.Context, opts []cehttp.Option) []cehttp.Option {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling opentelemetry middleware in http cloudevents server")

	otel.StartTracerProvider(ctx)
	otel.StartMeterProvider(ctx)

	optsotel := append(opts,
		ce.WithRoundTripper(otelhttp.NewTransport(http.DefaultTransport,
			otelhttp.WithTracerProvider(otel.TracerProvider),
			otelhttp.WithMeterProvider(otel.MeterProvider))),
		ce.WithMiddleware(func(next http.Handler) http.Handler {
			return otelhttp.NewHandler(next, "ce.http.receiver")
		}),
	)

	logger.Debug("opentelemetry integration successfully enabled in http cloudevents server")

	return optsotel
}
