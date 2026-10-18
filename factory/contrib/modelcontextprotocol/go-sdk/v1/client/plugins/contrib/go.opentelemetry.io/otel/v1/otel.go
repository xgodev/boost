// Package otel opens an OpenTelemetry span for every MCP call the client
// sends, propagating the trace to the server in _meta. Spans come from the
// boost otel factory's TracerProvider.
package otel

import (
	"context"

	bootel "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/internal/middleware"
	"go.opentelemetry.io/otel/trace"
)

// Scope is the instrumentation scope of the spans.
const Scope = "github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1"

// Register installs the tracing middleware on cli; a client.Plugin.
func Register(_ context.Context, cli *client.Client) error {
	o, err := NewOptions()
	if err != nil {
		return err
	}
	if !o.Enabled {
		return nil
	}
	cli.Instance().AddSendingMiddleware(middleware.Trace(bootel.TracerProvider.Tracer(Scope), trace.SpanKindClient))
	return nil
}
