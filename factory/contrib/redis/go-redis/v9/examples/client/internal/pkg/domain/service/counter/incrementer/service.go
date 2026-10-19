// Package incrementer increments counters.
package incrementer

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9/examples/client/internal/pkg/domain/model/entity"
)

// Service increments a counter.
type Service interface {
	Increment(ctx context.Context, name string) (*entity.Counter, error)
}
