// Package finder reads the units in stock of an item.
package finder

import (
	"context"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
)

// Service finds stock.
type Service interface {
	// Find returns the stock of sku; an unknown sku is a NotFound.
	Find(ctx context.Context, sku string) (*entity.Stock, error)
}
