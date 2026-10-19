// Package item is the persistence of the items.
package item

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/notfound/internal/pkg/domain/model/entity"
)

// Repository keeps the items.
type Repository interface {
	// Get returns the item of id; an unknown id is a NotFound.
	Get(ctx context.Context, id string) (*entity.Item, error)
}
