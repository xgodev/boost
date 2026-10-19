// Package stock is the persistence of the units in stock.
package stock

import (
	"context"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
)

// Repository keeps the units in stock of each item.
type Repository interface {
	// Create stores a new stock; a sku already stored is an AlreadyExists.
	Create(ctx context.Context, s *entity.Stock) error
	// Get returns the stock of sku; an unknown sku is a NotFound.
	Get(ctx context.Context, sku string) (*entity.Stock, error)
	// Take removes one unit of sku and returns what is left; an unknown sku
	// is a NotFound and a sold-out one a Conflict.
	Take(ctx context.Context, sku string) (*entity.Stock, error)
}
