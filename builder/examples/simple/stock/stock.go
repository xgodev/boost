// Package stock tells whether an item can be ordered. Its implementations,
// memory (the default), remote (asks a warehouse) and unlimited, are chosen
// by config: builder.examples.simple.stock.implementation.
package stock

import "context"

// Service reports the stock of items.
type Service interface {
	// Available reports whether item can be ordered.
	Available(ctx context.Context, item string) (bool, error)
}
