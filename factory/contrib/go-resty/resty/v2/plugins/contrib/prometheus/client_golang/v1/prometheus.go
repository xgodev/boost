package prometheus

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"errors"
	"github.com/go-resty/resty/v2"
	prom "github.com/prometheus/client_golang/prometheus"

	"github.com/xgodev/boost/wrapper/log"
)

var (
	reqTotal = prom.NewCounterVec(prom.CounterOpts{
		Name: "boost_factory_resty_requests_total",
		Help: "The number of requests made",
	}, []string{"code", "method", "host", "url"})

	reqDur = prom.NewHistogramVec(prom.HistogramOpts{
		Name:    "boost_factory_resty_request_duration_seconds",
		Help:    "The request latency in seconds",
		Buckets: prom.DefBuckets,
	}, []string{"code", "method", "host", "url"})
)

func init() {
	prom.MustRegister(reqDur)
	prom.MustRegister(reqTotal)
}

// Prometheus represents opentracing plugin for resty client.
type Prometheus struct {
	options *Options
}

// NewPrometheusWithConfigPath returns new opentracing with options from config path.
func NewPrometheusWithConfigPath(path string) (*Prometheus, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewPrometheusWithOptions(o), nil
}

// NewPrometheusWithOptions returns new opentracing with options.
func NewPrometheusWithOptions(options *Options) *Prometheus {
	return &Prometheus{options: options}
}

// Register registers a new opentracing plugin on resty client.
func Register(ctx context.Context, client *resty.Client) error {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	plugin := NewPrometheusWithOptions(o)
	return plugin.Register(ctx, client)
}

// Register registers this opentracing plugin on resty client.
func (p *Prometheus) Register(ctx context.Context, client *resty.Client) error {

	if !p.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)
	logger.Trace("enabling prometheus middleware in resty")

	client.OnBeforeRequest(p.beforeRequest).
		OnAfterResponse(p.collectAfterResponse).
		OnError(p.collectError)

	logger.Debug("prometheus middleware successfully enabled in resty")

	return nil
}

const defaultSubsystem = "resty"

type ctxkey struct{ name string }

var urlKey = &ctxkey{"request-url"}

func (p *Prometheus) collect(req *http.Request, code int, dur time.Duration) {
	url, _ := req.Context().Value(urlKey).(string)

	values := []string{
		strconv.Itoa(code),
		req.Method,
		req.URL.Hostname(),
		url,
	}

	reqTotal.WithLabelValues(values...).Inc()
	reqDur.WithLabelValues(values...).Observe(dur.Seconds())
}

func (p *Prometheus) beforeRequest(client *resty.Client, req *resty.Request) error {
	ctx := context.WithValue(req.Context(), urlKey, req.URL)
	req.SetContext(ctx)

	return nil
}

func (p *Prometheus) collectAfterResponse(client *resty.Client, res *resty.Response) error {
	p.collect(
		res.Request.RawRequest,
		res.StatusCode(),
		res.Time(),
	)

	return nil
}

func (p *Prometheus) collectError(req *resty.Request, err error) {
	code := http.StatusInternalServerError

	var dur time.Duration
	var e *resty.ResponseError

	if errors.As(err, &e) {
		code = e.Response.StatusCode()
		dur = e.Response.Time()
	}

	p.collect(
		req.RawRequest,
		code,
		dur,
	)
}
