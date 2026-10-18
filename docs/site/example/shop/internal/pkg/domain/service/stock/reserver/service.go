// Package reserver reserves one unit of an item.
package reserver

import "context"

// Service reserves stock.
type Service interface {
	Reserve(ctx context.Context, sku string) error
}
