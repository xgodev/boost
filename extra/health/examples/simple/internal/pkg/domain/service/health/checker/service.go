// Package checker checks the health of the app.
package checker

import (
	"context"

	"github.com/xgodev/boost/extra/health/examples/simple/internal/pkg/domain/model/entity"
)

// Service checks the health of the app.
type Service interface {
	// Check runs every health checker.
	Check(ctx context.Context) []entity.Check
}
