// Package stock is the persistence of the units in stock.
package stock

import "context"

// Repository keeps the units in stock of each item.
type Repository interface {
	// Take removes one unit of sku; a sold-out sku is a Conflict.
	Take(ctx context.Context, sku string) error
}
