// Package creator puts a new item in stock.
package creator

import (
	"context"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
)

// Service creates stock.
type Service interface {
	// Create stores the stock of a new sku; a blank sku or negative units is
	// a NotValid and a sku already in stock an AlreadyExists.
	Create(ctx context.Context, s *entity.Stock) error
}
