// Package note is the persistence of the notes.
package note

import (
	"context"

	"github.com/xgodev/boost/wrapper/cache/examples/plugin/internal/pkg/domain/model/entity"
)

// Repository keeps the notes.
type Repository interface {
	Save(ctx context.Context, note *entity.Note) error
	Get(ctx context.Context, key string) (*entity.Note, error)
}
