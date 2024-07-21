package otel

import (
	"context"
	"go.opentelemetry.io/otel/propagation"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/credentials"
)

var tracerProvider trace.TracerProvider

// StartTracerProvider starts the tracer provider like StartMetricProviderWithOptions but with default Options.
func StartTracerProvider(ctx context.Context, startOptions ...sdktrace.TracerProviderOption) {

	o, err := NewOptions()
	if err != nil {
		panic(err)
	}

	StartTracerProviderWithOptions(ctx, o, startOptions...)
}

var tracerOnce sync.Once

// StartTracerProviderWithOptions starts the tracer provider with the given set of options. Calling
// it multiple times will have no effect. If an error occours during tracer initialization,
// a Noop trace provider will be used instead.
func StartTracerProviderWithOptions(ctx context.Context, options *Options, startOptions ...sdktrace.TracerProviderOption) {

	if !IsTraceEnabled() {
		return
	}

	tracerOnce.Do(func() {

		tracerProvider = noop.NewTracerProvider()

		logger := log.FromContext(ctx)

		exporter, err := NewTracerExporter(ctx, options)

		if err != nil {
			logger.Error("error creating opentelemetry exporter: ", err)
			otel.SetTracerProvider(noop.NewTracerProvider())
			return
		}

		rs, err := NewResource(ctx, options)
		if err != nil {
			logger.Error("error creating opentelemetry resource: ", err)
			otel.SetTracerProvider(noop.NewTracerProvider())
			return
		}

		startOptions = append(startOptions,
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(rs),
		)

		prov := sdktrace.NewTracerProvider(startOptions...)

		otel.SetTracerProvider(prov)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		tracerProvider = prov

		log.Infof("started opentelemetry tracer: %s", options.Service)
	})
}

func NewTracerExporter(ctx context.Context, options *Options) (*otlptrace.Exporter, error) {
	var exporter *otlptrace.Exporter
	var err error

	switch options.Protocol {
	case "grpc":
		exporter, err = NewGRPCTracerExporter(ctx, options)
	default:
		exporter, err = NewHTTPTracerExporter(ctx, options)
	}
	return exporter, err
}

func NewHTTPTracerExporter(ctx context.Context, options *Options) (*otlptrace.Exporter, error) {
	var exporterOpts []otlptracehttp.Option
	if _, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT"); !ok { // Only using WithEndpoint when the environment variable is not set
		exporterOpts = append(exporterOpts, otlptracehttp.WithEndpoint(options.Endpoint)) //TODO see https://github.com/open-telemetry/opentelemetry-go/issues/3730
	}

	if IsInsecure() {
		exporterOpts = append(exporterOpts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(
		ctx,
		exporterOpts...,
	)
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

func NewGRPCTracerExporter(ctx context.Context, options *Options) (*otlptrace.Exporter, error) {
	var exporterOpts []otlptracegrpc.Option
	if _, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT"); !ok { // Only using WithEndpoint when the environment variable is not set
		exporterOpts = append(exporterOpts, otlptracegrpc.WithEndpoint(options.Endpoint)) //TODO see https://github.com/open-telemetry/opentelemetry-go/issues/3730
	}

	if IsInsecure() {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithInsecure())
	} else {
		creds, err := credentials.NewClientTLSFromFile(options.TLS.Cert, "")
		if err != nil {
			return nil, errors.Wrap(err, "error creating tls credentials")
		}
		exporterOpts = append(exporterOpts, otlptracegrpc.WithTLSCredentials(creds))
	}

	exporter, err := otlptracegrpc.New(
		ctx,
		exporterOpts...,
	)
	if err != nil {
		return nil, err
	}

	return exporter, nil
}

// NewTracer creates a Tracer with the provided name and options. A Tracer
// allows the creation of spans for custom instrumentation.
//
// StartMetricProvider should be called before to setup the tracer provider, otherwise a Noop
// tracer provider will be used.
func NewTracer(name string, options ...trace.TracerOption) trace.Tracer {
	return tracerProvider.Tracer(name, options...)
}
