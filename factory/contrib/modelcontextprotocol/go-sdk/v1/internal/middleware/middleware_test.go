package middleware

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type entry struct {
	level  string
	msg    string
	fields map[string]any
}

// recorder is a log.Logger that keeps what the middleware writes.
type recorder struct {
	log.Noop
	mu      *sync.Mutex
	entries *[]entry
	fields  map[string]any
}

func newRecorder() recorder {
	return recorder{mu: new(sync.Mutex), entries: new([]entry)}
}

func (r recorder) FromContext(context.Context) log.Logger { return r }

func (r recorder) WithField(k string, v any) log.Logger {
	f := maps.Clone(r.fields)
	if f == nil {
		f = map[string]any{}
	}
	f[k] = v
	r.fields = f
	return r
}

func (r recorder) add(level, format string, args ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	*r.entries = append(*r.entries, entry{level: level, msg: fmt.Sprintf(format, args...), fields: r.fields})
}

func (r recorder) Tracef(f string, a ...any) { r.add("TRACE", f, a...) }
func (r recorder) Debugf(f string, a ...any) { r.add("DEBUG", f, a...) }
func (r recorder) Infof(f string, a ...any)  { r.add("INFO", f, a...) }
func (r recorder) Errorf(f string, a ...any) { r.add("ERROR", f, a...) }

func (r recorder) find(method string) []entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []entry
	for _, e := range *r.entries {
		if e.fields[FieldMethod] == method {
			out = append(out, e)
		}
	}
	return out
}

func useRecorder(t *testing.T) recorder {
	t.Helper()
	r := newRecorder()
	log.Set(r)
	t.Cleanup(func() { log.Set(log.Noop{}) })
	return r
}

type echoIn struct {
	Text string `json:"text"`
}

func pair(t *testing.T, serverMW, clientMW []mcp.Middleware) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()

	srv := mcp.NewServer(&mcp.Implementation{Name: "srv", Version: "1.0.0"}, nil)
	mcp.AddTool(srv, &mcp.Tool{Name: "echo"},
		func(_ context.Context, _ *mcp.CallToolRequest, in echoIn) (*mcp.CallToolResult, any, error) {
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: in.Text}}}, nil, nil
		})
	mcp.AddTool(srv, &mcp.Tool{Name: "fail"},
		func(context.Context, *mcp.CallToolRequest, echoIn) (*mcp.CallToolResult, any, error) {
			return nil, nil, fmt.Errorf("tool broke")
		})
	srv.AddReceivingMiddleware(serverMW...)

	st, ct := mcp.NewInMemoryTransports()
	ss, err := srv.Connect(ctx, st, nil)
	require.NoError(t, err)
	t.Cleanup(func() { ss.Close() })

	cli := mcp.NewClient(&mcp.Implementation{Name: "cli", Version: "1.0.0"}, nil)
	cli.AddSendingMiddleware(clientMW...)
	cs, err := cli.Connect(ctx, ct, nil)
	require.NoError(t, err)
	t.Cleanup(func() { cs.Close() })
	return cs
}

func call(cs *mcp.ClientSession, tool string) (*mcp.CallToolResult, error) {
	return cs.CallTool(context.Background(), &mcp.CallToolParams{Name: tool, Arguments: map[string]any{"text": "hi"}})
}

func TestLog_LogsMethodToolAndLevel(t *testing.T) {
	rec := useRecorder(t)
	cs := pair(t, []mcp.Middleware{Log("INFO")}, []mcp.Middleware{Log("TRACE")})

	_, err := call(cs, "echo")
	require.NoError(t, err)

	got := rec.find("tools/call")
	require.Len(t, got, 2, "one line per side")
	assert.ElementsMatch(t, []string{"INFO", "TRACE"}, []string{got[0].level, got[1].level})
	for _, e := range got {
		assert.Equal(t, "echo", e.fields[FieldTool])
		assert.Contains(t, e.msg, "tools/call echo")
	}
}

func TestLog_DefaultLevelIsDebug(t *testing.T) {
	rec := useRecorder(t)
	cs := pair(t, nil, []mcp.Middleware{Log("")})

	_, err := call(cs, "echo")
	require.NoError(t, err)

	got := rec.find("tools/call")
	require.Len(t, got, 1)
	assert.Equal(t, "DEBUG", got[0].level)
}

func TestLog_ErrorsGoToErrorAndToolErrorsAreFlagged(t *testing.T) {
	rec := useRecorder(t)
	cs := pair(t, []mcp.Middleware{Log("INFO")}, nil)

	_, err := call(cs, "nope")
	assert.Error(t, err)
	res, err := call(cs, "fail")
	require.NoError(t, err)
	require.True(t, res.IsError)

	got := rec.find("tools/call")
	require.Len(t, got, 2)
	assert.Equal(t, "ERROR", got[0].level)
	assert.Contains(t, got[0].msg, "nope")
	assert.Equal(t, "INFO", got[1].level)
	assert.Contains(t, got[1].msg, "isError")
}

func TestLog_SilencedErrorDropsToDebug(t *testing.T) {
	rec := useRecorder(t)
	silenced := fmt.Errorf("client went away")
	errors.Ignore(silenced, errors.IgnoreSilenceLog)
	t.Cleanup(errors.ResetRegistry)

	h := Log("INFO")(func(context.Context, string, mcp.Request) (mcp.Result, error) {
		return nil, silenced
	})
	_, err := h(context.Background(), "tools/call", &mcp.CallToolRequest{Params: &mcp.CallToolParamsRaw{Name: "x"}})
	assert.ErrorIs(t, err, silenced)

	got := rec.find("tools/call")
	require.Len(t, got, 1)
	assert.Equal(t, "DEBUG", got[0].level)
}

func TestLog_NilParamsHaveNoTool(t *testing.T) {
	rec := useRecorder(t)

	h := Log("INFO")(func(context.Context, string, mcp.Request) (mcp.Result, error) { return nil, nil })
	_, err := h(context.Background(), "tools/list", &mcp.ListToolsRequest{})
	require.NoError(t, err)

	got := rec.find("tools/list")
	require.Len(t, got, 1)
	assert.NotContains(t, got[0].fields, FieldTool)
}

func tracing(t *testing.T) (*tracetest.SpanRecorder, trace.Tracer) {
	t.Helper()
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	prev := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() { otel.SetTextMapPropagator(prev) })
	return sr, tp.Tracer("test")
}

func spansNamed(sr *tracetest.SpanRecorder, name string, kind trace.SpanKind) []sdktrace.ReadOnlySpan {
	var out []sdktrace.ReadOnlySpan
	for _, s := range sr.Ended() {
		if s.Name() == name && s.SpanKind() == kind {
			out = append(out, s)
		}
	}
	return out
}

func attr(s sdktrace.ReadOnlySpan, key string) string {
	for _, a := range s.Attributes() {
		if string(a.Key) == key {
			return a.Value.Emit()
		}
	}
	return ""
}

func TestTrace_PropagatesThroughMetaAndNamesSpans(t *testing.T) {
	sr, tracer := tracing(t)
	cs := pair(t,
		[]mcp.Middleware{Trace(tracer, trace.SpanKindServer)},
		[]mcp.Middleware{Trace(tracer, trace.SpanKindClient)})

	_, err := call(cs, "echo")
	require.NoError(t, err)

	client := spansNamed(sr, "tools/call echo", trace.SpanKindClient)
	server := spansNamed(sr, "tools/call echo", trace.SpanKindServer)
	require.Len(t, client, 1)
	require.Len(t, server, 1)

	assert.Equal(t, client[0].SpanContext().TraceID(), server[0].SpanContext().TraceID())
	assert.Equal(t, client[0].SpanContext().SpanID(), server[0].Parent().SpanID())
	assert.Equal(t, "tools/call", attr(server[0], AttrMethod))
	assert.Equal(t, "echo", attr(server[0], AttrTool))
	assert.Equal(t, codes.Unset, server[0].Status().Code)
}

func TestTrace_DoesNotMutateCallerMeta(t *testing.T) {
	_, tracer := tracing(t)
	cs := pair(t, nil, []mcp.Middleware{Trace(tracer, trace.SpanKindClient)})

	meta := mcp.Meta{"tenant": "a"}
	_, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Meta: meta, Name: "echo", Arguments: map[string]any{"text": "hi"}})
	require.NoError(t, err)
	// The SDK adds its own keys to the caller's map; the trace context must not land there.
	assert.NotContains(t, meta, "traceparent")
	assert.Equal(t, "a", meta["tenant"])
}

func TestTrace_MarksErrorsAndToolErrors(t *testing.T) {
	sr, tracer := tracing(t)
	cs := pair(t, []mcp.Middleware{Trace(tracer, trace.SpanKindServer)}, nil)

	_, err := call(cs, "nope")
	assert.Error(t, err)
	_, err = call(cs, "fail")
	require.NoError(t, err)

	nope := spansNamed(sr, "tools/call nope", trace.SpanKindServer)
	require.Len(t, nope, 1)
	assert.Equal(t, codes.Error, nope[0].Status().Code)
	assert.NotEmpty(t, nope[0].Events(), "error recorded on the span")

	fail := spansNamed(sr, "tools/call fail", trace.SpanKindServer)
	require.Len(t, fail, 1)
	assert.Equal(t, codes.Error, fail[0].Status().Code)
	assert.Equal(t, "tool_error", attr(fail[0], AttrErrorType))
}

func TestTrace_NilParamsAndSession(t *testing.T) {
	sr, tracer := tracing(t)

	for _, kind := range []trace.SpanKind{trace.SpanKindServer, trace.SpanKindClient} {
		h := Trace(tracer, kind)(func(context.Context, string, mcp.Request) (mcp.Result, error) { return nil, nil })
		_, err := h(context.Background(), "tools/list", &mcp.ListToolsRequest{})
		require.NoError(t, err)
	}
	require.Len(t, sr.Ended(), 2)
	assert.Empty(t, attr(sr.Ended()[0], AttrSession))
}

func TestTrace_RecordsSessionIDOnStreamableHTTP(t *testing.T) {
	sr, tracer := tracing(t)
	srv := mcp.NewServer(&mcp.Implementation{Name: "srv", Version: "1.0.0"}, nil)
	srv.AddReceivingMiddleware(Trace(tracer, trace.SpanKindServer))
	ts := httptest.NewServer(mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return srv }, nil))
	defer ts.Close()

	cli := mcp.NewClient(&mcp.Implementation{Name: "cli", Version: "1.0.0"}, nil)
	cs, err := cli.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: ts.URL, MaxRetries: -1}, nil)
	require.NoError(t, err)
	defer cs.Close()
	_, err = cs.ListTools(context.Background(), nil)
	require.NoError(t, err)

	spans := spansNamed(sr, "tools/list", trace.SpanKindServer)
	require.Len(t, spans, 1)
	assert.Equal(t, cs.ID(), attr(spans[0], AttrSession))
}

func TestMetaCarrier(t *testing.T) {
	c := metaCarrier{"traceparent": "00-x", "count": 1}
	assert.Equal(t, "00-x", c.Get("traceparent"))
	assert.Empty(t, c.Get("count"), "non-string values read as empty")
	assert.ElementsMatch(t, []string{"traceparent", "count"}, c.Keys())
	assert.True(t, isNil(nil))
}

func TestTrace_ClientInjectsIntoEmptyMeta(t *testing.T) {
	_, tracer := tracing(t)
	params := &mcp.CallToolParamsRaw{Name: "x"}

	h := Trace(tracer, trace.SpanKindClient)(func(context.Context, string, mcp.Request) (mcp.Result, error) { return nil, nil })
	_, err := h(context.Background(), "tools/call", &mcp.CallToolRequest{Params: params})
	require.NoError(t, err)
	assert.NotEmpty(t, params.Meta["traceparent"])
}
