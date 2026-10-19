// Package greeter greets.
package greeter

import (
	"context"

	"github.com/xgodev/boost/factory/local/extra/middleware/examples/error/internal/pkg/domain/model/entity"
)

// Service greets.
type Service interface {
	// Greet greets name in language lang.
	Greet(ctx context.Context, name, lang string) (*entity.Greeting, error)
}
