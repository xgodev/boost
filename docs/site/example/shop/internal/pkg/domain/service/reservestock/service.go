// Package reservestock reserves one unit of an item.
package reservestock

import "context"

// Service reserves stock.
type Service interface {
	Reserve(ctx context.Context, sku string) error
}
