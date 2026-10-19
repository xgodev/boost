// Package keeper keeps the readings of the sensors.
package keeper

import (
	"context"

	"github.com/xgodev/boost/wrapper/cache/examples/freecache-binary/internal/pkg/domain/model/entity"
)

// Service keeps the last reading of each sensor.
type Service interface {
	Save(ctx context.Context, r *entity.Reading) error
	Find(ctx context.Context, sensor uint32) (*entity.Reading, error)
}
