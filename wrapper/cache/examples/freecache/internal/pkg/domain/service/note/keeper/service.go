// Package keeper keeps notes.
package keeper

import (
	"context"

	"github.com/xgodev/boost/wrapper/cache/examples/freecache/internal/pkg/domain/model/entity"
)

// Service keeps notes.
type Service interface {
	Save(ctx context.Context, note *entity.Note) error
	Find(ctx context.Context, key string) (*entity.Note, error)
}
