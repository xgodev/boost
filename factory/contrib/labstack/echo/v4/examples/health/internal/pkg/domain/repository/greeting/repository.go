// Package greeting is the persistence of the greeting.
package greeting

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/health/internal/pkg/domain/model/entity"
)

// Repository keeps the greeting.
type Repository interface {
	Get(ctx context.Context) (*entity.Greeting, error)
}
