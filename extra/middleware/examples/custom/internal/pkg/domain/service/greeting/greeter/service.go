// Package greeter greets.
package greeter

import (
	"context"

	"github.com/xgodev/boost/extra/middleware/examples/custom/internal/pkg/domain/model/entity"
)

// Service greets.
type Service interface {
	Greet(ctx context.Context, name string) (*entity.Greeting, error)
}
