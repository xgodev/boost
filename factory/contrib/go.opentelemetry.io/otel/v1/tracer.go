package otel

import (
	"context"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/propagation"
	"sync"

	"github.com/pkg/errors"
	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/credentials"
)

var TracerProvider trace.TracerProvider

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

		TracerProvider = noop.NewTracerProvider()

		logger := log.FromContext(ctx)

		otel.SetLogger(logr.New(&Logger{}))

		exporter, err := NewTracerExporter(ctx, options)

		if err != nil {
			logger.Error("error creating opentelemetry exporter: ", err)
			return
		}

		rs, err := NewResource(ctx, options)
		if err != nil {
			logger.Error("error creating opentelemetry resource: ", err)
			return
		}

		if options.Console.Enabled {
			// 1) Console exporter
			consoleExp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
			if err != nil {
				logger.Error("stdout trace exporter: %v", err)
			}
			startOptions = append(startOptions, sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(consoleExp)))
		}

		startOptions = append(startOptions,
			sdktrace.WithBatcher(exporter),
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(rs),
		)

		prov := sdktrace.NewTracerProvider(startOptions...)

		otel.SetTracerProvider(prov)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		TracerProvider = prov

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

	exporterOpts = append(exporterOpts, otlptracehttp.WithEndpoint(options.Endpoint))

	if IsInsecure() {
		exporterOpts = append(exporterOpts, otlptracehttp.WithInsecure())
	}

	return otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			exporterOpts...,
		),
	)
}

func NewGRPCTracerExporter(ctx context.Context, options *Options) (*otlptrace.Exporter, error) {
	exporterOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(options.Endpoint),
	}

	if IsInsecure() {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithInsecure())
	} else if options.TLS.Cert != "" {
		creds, err := credentials.NewClientTLSFromFile(options.TLS.Cert, "")
		if err != nil {
			return nil, errors.Wrap(err, "error creating tls credentials")
		}
		exporterOpts = append(exporterOpts, otlptracegrpc.WithTLSCredentials(creds))
	}

	return otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			exporterOpts...,
		),
	)
}

// NewTracer creates a Tracer with the provided name and options. A Tracer
// allows the creation of spans for custom instrumentation.
//
// StartTracerProvider should be called before to setup the tracer provider, otherwise a Noop
// tracer provider will be used.
func NewTracer(name string, options ...trace.TracerOption) trace.Tracer {
	return TracerProvider.Tracer(name, options...)
}
