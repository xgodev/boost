// Package reserver reserves one unit of an item. Its implementations, local
// (the default) and resty (asks app api over HTTP), are chosen by
// internal.pkg.domain.service.stock.reserver.implementation.
package reserver

import (
	"context"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
)

// Service reserves stock.
type Service interface {
	// Reserve takes one unit of sku and returns what is left; an unknown
	// sku is a NotFound and a sold-out one a Conflict.
	Reserve(ctx context.Context, sku string) (*entity.Stock, error)
}
