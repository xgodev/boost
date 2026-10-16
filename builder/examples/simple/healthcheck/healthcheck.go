// Package healthcheck serves grpc.health.v1.Health on the grpc server of
// app api.
package healthcheck

import (
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// NewServer provides the health service, registered by the generated grpc
// door with healthpb.RegisterHealthServer.
// @Provide(index=0)
// @GRPC(register=healthpb.RegisterHealthServer)
func NewServer() *health.Server {
	srv := health.NewServer()
	srv.SetServingStatus("orders", healthpb.HealthCheckResponse_SERVING)
	return srv
}
