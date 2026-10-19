// Package counter keeps the counters.
package counter

import "context"

// Repository keeps the counters.
type Repository interface {
	// Increment adds one to the counter and returns its count.
	Increment(ctx context.Context, name string) (int64, error)
}
