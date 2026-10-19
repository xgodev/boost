// Package finder finds an item.
package finder

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/notfound/internal/pkg/domain/model/entity"
)

// Service finds items.
type Service interface {
	// Find returns the item of id; an unknown id is a NotFound.
	Find(ctx context.Context, id string) (*entity.Item, error)
}
