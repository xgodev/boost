package otel

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewResource(ctx context.Context, options *Options) (*resource.Resource, error) {

	attrs := make([]attribute.KeyValue, len(options.Tags))
	for k, v := range options.Tags {
		attrs = append(attrs, attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(v),
		})
	}

	return resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(options.Service),
			semconv.ServiceVersionKey.String(options.Version),
			attribute.String("env", options.Env),
		),
		resource.WithAttributes(attrs...),
	)

}
