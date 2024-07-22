package otel

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func NewResource(ctx context.Context, options *Options) (*resource.Resource, error) {

	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(options.Service),
		semconv.ServiceVersionKey.String(options.Version),
		attribute.String("env", options.Env),
		attribute.String("library.language", "go"),
	}

	for k, v := range options.Tags {
		attrs = append(attrs, attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(v),
		})
	}

	return resource.New(ctx, resource.WithAttributes(
		attrs...,
	))

}
