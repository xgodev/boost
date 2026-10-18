// Package middleware holds the go-sdk middlewares the server and client
// plugins install: one call in, one log line or span out.
package middleware

import (
	"context"
	"maps"
	"reflect"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Log fields and span attributes. The span names follow the OpenTelemetry
// MCP semantic conventions.
const (
	FieldMethod = "mcp.method"
	FieldTool   = "mcp.tool"

	AttrMethod    = "mcp.method.name"
	AttrTool      = "gen_ai.tool.name"
	AttrSession   = "mcp.session.id"
	AttrErrorType = "error.type"
)

// Log logs every call at level (TRACE, INFO, otherwise DEBUG). Errors go to
// Error unless errors.IgnoreOf silences them; a tool result with isError is
// logged at level and flagged.
func Log(level string) mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			start := time.Now()
			res, err := next(ctx, method, req)
			elapsed := time.Since(start)

			logger := log.FromContext(ctx).WithField(FieldMethod, method)
			target := method
			if tool := toolName(req); tool != "" {
				logger = logger.WithField(FieldTool, tool)
				target += " " + tool
			}

			if err != nil {
				if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreSilenceLog != 0 {
					logger.Debugf("mcp %s %v: %v", target, elapsed, err)
				} else {
					logger.Errorf("mcp %s %v: %v", target, elapsed, err)
				}
				return res, err
			}

			logf := logger.Debugf
			switch level {
			case "TRACE":
				logf = logger.Tracef
			case "INFO":
				logf = logger.Infof
			}
			if isToolError(res) {
				logf("mcp %s %v isError", target, elapsed)
			} else {
				logf("mcp %s %v", target, elapsed)
			}
			return res, err
		}
	}
}

// Trace opens a span per call. SpanKindServer continues the trace the caller
// put in params._meta; SpanKindClient writes the current one there.
func Trace(tracer trace.Tracer, kind trace.SpanKind) mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			params := req.GetParams()
			if isNil(params) {
				params = nil
			}

			attrs := []attribute.KeyValue{attribute.String(AttrMethod, method)}
			name := method
			if tool := toolName(req); tool != "" {
				attrs = append(attrs, attribute.String(AttrTool, tool))
				name += " " + tool
			}
			if s := req.GetSession(); !isNil(s) && s.ID() != "" {
				attrs = append(attrs, attribute.String(AttrSession, s.ID()))
			}

			if kind == trace.SpanKindServer && params != nil {
				ctx = otel.GetTextMapPropagator().Extract(ctx, metaCarrier(params.GetMeta()))
			}

			ctx, span := tracer.Start(ctx, name, trace.WithSpanKind(kind), trace.WithAttributes(attrs...))
			defer span.End()

			if kind == trace.SpanKindClient && params != nil {
				// A copy: the caller's Meta map is theirs and may be reused.
				meta := maps.Clone(params.GetMeta())
				if meta == nil {
					meta = map[string]any{}
				}
				otel.GetTextMapPropagator().Inject(ctx, metaCarrier(meta))
				params.SetMeta(meta)
			}

			res, err := next(ctx, method, req)
			switch {
			case err != nil:
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			case isToolError(res):
				span.SetAttributes(attribute.String(AttrErrorType, "tool_error"))
				span.SetStatus(codes.Error, "tool returned isError")
			}
			return res, err
		}
	}
}

func toolName(req mcp.Request) string {
	switch p := req.GetParams().(type) {
	case *mcp.CallToolParams:
		if p != nil {
			return p.Name
		}
	case *mcp.CallToolParamsRaw:
		if p != nil {
			return p.Name
		}
	}
	return ""
}

func isToolError(res mcp.Result) bool {
	r, ok := res.(*mcp.CallToolResult)
	return ok && r != nil && r.IsError
}

// isNil reports a nil interface or one holding a nil pointer: the SDK passes
// typed nils (a ping has *PingParams(nil), a bare request a nil *ServerSession).
func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Pointer && rv.IsNil()
}

// metaCarrier lets the OTel propagator read and write params._meta.
type metaCarrier map[string]any

func (c metaCarrier) Get(key string) string {
	s, _ := c[key].(string)
	return s
}

func (c metaCarrier) Set(key, value string) { c[key] = value }

func (c metaCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
