// Package describer describes the application.
package describer

import (
	"context"

	"github.com/xgodev/boost/wrapper/config/examples/env/internal/pkg/domain/model/entity"
)

// Service describes the application.
type Service interface {
	Describe(ctx context.Context) (*entity.Application, error)
}
