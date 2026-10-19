// Package local checks the health of the app with boost's health registry.
package local

import (
	"context"

	"github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/extra/health/examples/simple/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/extra/health/examples/simple/internal/pkg/domain/service/health/checker"
	"github.com/xgodev/boost/extra/health/examples/simple/internal/pkg/health/tcp"
)

type service struct{}

// New adds the tcp checker to boost's health registry, required, and
// provides the service that runs the registry.
// @Provide(index=0)
func New(dependency *tcp.Checker) checker.Service {
	health.Add(health.NewHealthChecker("dependency", "the tcp dependency", dependency, true, true))
	return &service{}
}

func (s *service) Check(ctx context.Context) []entity.Check {
	var checks []entity.Check
	for _, result := range health.CheckAll(ctx) {
		check := entity.Check{Name: result.HealthCheck.Name, Required: result.HealthCheck.IsRequired()}
		if !result.IsOk() {
			check.Error = result.Error.Error()
		}
		checks = append(checks, check)
	}
	return checks
}
