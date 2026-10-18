// Package otel opens an OpenTelemetry span for every MCP call the server
// receives, continuing the trace the client sent in _meta. Spans come from the
// boost otel factory's TracerProvider.
package otel

import (
	"context"

	bootel "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/internal/middleware"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	"go.opentelemetry.io/otel/trace"
)

// Scope is the instrumentation scope of the spans.
const Scope = "github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1"

// Register installs the tracing middleware on srv; a server.Plugin.
func Register(_ context.Context, srv *server.Server) error {
	o, err := NewOptions()
	if err != nil {
		return err
	}
	if !o.Enabled {
		return nil
	}
	srv.Instance().AddReceivingMiddleware(middleware.Trace(bootel.TracerProvider.Tracer(Scope), trace.SpanKindServer))
	return nil
}
