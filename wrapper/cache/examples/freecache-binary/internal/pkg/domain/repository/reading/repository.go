// Package reading is the persistence of the readings.
package reading

import (
	"context"

	"github.com/xgodev/boost/wrapper/cache/examples/freecache-binary/internal/pkg/domain/model/entity"
)

// Repository keeps the last reading of each sensor.
type Repository interface {
	Save(ctx context.Context, r *entity.Reading) error
	Get(ctx context.Context, sensor uint32) (*entity.Reading, error)
}
