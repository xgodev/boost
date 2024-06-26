package datadog

import (
	"context"
	"strconv"

	"github.com/go-resty/resty/v2"
	datadog "github.com/xgodev/boost/factory/contrib/datadog/dd-trace-go/v1"
	"github.com/xgodev/boost/wrapper/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Datadog represents Datadog integration with resty.
type Datadog struct {
	options *Options
}

// NewDatadogWithConfigPath returns a new Datadog with options from config path.
func NewDatadogWithConfigPath(path string, spanOptions ...ddtrace.StartSpanOption) (*Datadog, error) {
	o, err := NewOptionsWithPath(path, spanOptions...)
	if err != nil {
		return nil, err
	}
	return NewDatadogWithOptions(o), nil
}

// NewDatadog returns a new DataDog with default options.
func NewDatadogWithOptions(options *Options) *Datadog {
	return &Datadog{options: options}
}

// NewDatadog returns a new DataDog with default options.
func NewDatadog(traceOptions ...ddtrace.StartSpanOption) *Datadog {
	o, err := NewOptions(traceOptions...)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return NewDatadogWithOptions(o)
}

// Register registers Datadog integration with resty. It is shorthand for NewDatadog().Register(ctx, client).
func Register(ctx context.Context, client *resty.Client) error {
	o, err := NewOptions()
	if err != nil {
		return err
	}
	d := NewDatadogWithOptions(o)
	return d.Register(ctx, client)
}

// Register registers Datadog integration with resty.
func (d *Datadog) Register(ctx context.Context, client *resty.Client) error {
	if !d.options.Enabled || !datadog.IsTracerEnabled() {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("integrating resty in datadog")

	client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		spanOptions := []ddtrace.StartSpanOption{
			tracer.ResourceName(request.URL),
			tracer.SpanType(ext.SpanTypeHTTP),
			tracer.Tag(ext.HTTPMethod, request.Method),
			tracer.Tag(ext.HTTPURL, request.URL),
		}

		spanOptions = append(spanOptions, d.options.SpanOptions...)

		reqCtx := request.Context()
		span, ctx := tracer.StartSpanFromContext(reqCtx, d.options.OperationName, spanOptions...)

		// pass the span through the request context
		request.SetContext(ctx)

		return tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(request.Header))
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		ctx := resp.Request.Context()

		span, ok := tracer.SpanFromContext(ctx)

		if ok {
			span.SetTag(ext.HTTPCode, strconv.Itoa(resp.StatusCode()))
			span.SetTag(ext.Error, resp.Error())
			span.Finish()
		}

		return nil
	})

	logger.Debug("resty successfully integrated in datadog")

	return nil
}
