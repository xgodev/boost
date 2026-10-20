// Package oteltest runs an in-process OTLP/gRPC collector, so a test can
// point boost's otel factory at it (BOOST_FACTORY_OTEL_ENDPOINT, set before
// boost.Start) and assert which spans and metrics were exported, with no
// collector running on the machine.
package oteltest

import (
	"context"
	"net"
	"sync"

	colmetricpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
)

// Collector receives OTLP exports over gRPC on a loopback port and keeps the
// names of the spans and metrics it received.
type Collector struct {
	server   *grpc.Server
	endpoint string

	mu      sync.Mutex
	spans   []string
	metrics []string
}

// NewCollector starts a Collector on a free loopback port. Close stops it.
func NewCollector() (*Collector, error) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	c := &Collector{server: grpc.NewServer(), endpoint: lis.Addr().String()}
	coltracepb.RegisterTraceServiceServer(c.server, traceService{c: c})
	colmetricpb.RegisterMetricsServiceServer(c.server, metricService{c: c})

	go func() {
		_ = c.server.Serve(lis) // intentional: Serve only returns once Close stops the server
	}()

	return c, nil
}

// Endpoint is the host:port the collector listens on, the value for
// BOOST_FACTORY_OTEL_ENDPOINT (the default protocol is grpc, insecure).
func (c *Collector) Endpoint() string {
	return c.endpoint
}

// Spans returns the names of the spans received so far, in arrival order.
func (c *Collector) Spans() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]string(nil), c.spans...)
}

// Metrics returns the names of the metrics received so far, in arrival order.
func (c *Collector) Metrics() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]string(nil), c.metrics...)
}

// Close stops the collector and its open connections.
func (c *Collector) Close() {
	c.server.Stop()
}

type traceService struct {
	coltracepb.UnimplementedTraceServiceServer
	c *Collector
}

func (s traceService) Export(_ context.Context, req *coltracepb.ExportTraceServiceRequest) (*coltracepb.ExportTraceServiceResponse, error) {
	s.c.mu.Lock()
	defer s.c.mu.Unlock()
	for _, rs := range req.GetResourceSpans() {
		for _, ss := range rs.GetScopeSpans() {
			for _, span := range ss.GetSpans() {
				s.c.spans = append(s.c.spans, span.GetName())
			}
		}
	}
	return &coltracepb.ExportTraceServiceResponse{}, nil
}

type metricService struct {
	colmetricpb.UnimplementedMetricsServiceServer
	c *Collector
}

func (s metricService) Export(_ context.Context, req *colmetricpb.ExportMetricsServiceRequest) (*colmetricpb.ExportMetricsServiceResponse, error) {
	s.c.mu.Lock()
	defer s.c.mu.Unlock()
	for _, rm := range req.GetResourceMetrics() {
		for _, sm := range rm.GetScopeMetrics() {
			for _, m := range sm.GetMetrics() {
				s.c.metrics = append(s.c.metrics, m.GetName())
			}
		}
	}
	return &colmetricpb.ExportMetricsServiceResponse{}, nil
}
