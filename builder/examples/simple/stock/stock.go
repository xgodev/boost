// Package stock tells whether an item can be ordered. Its methods are plain:
// builder runs the middlewares @Middleware lists around each of them.
package stock

import "context"

// Service reports the stock of items.
type Service interface {
	// Available reports whether item can be ordered.
	Available(ctx context.Context, item string) (bool, error)
}

type memory struct {
	soldOut map[string]bool
}

// New provides the stock service; every consumer receives it wrapped in the
// log and prometheus middlewares.
// @Provide(index=0)
// @Middleware(log, prometheus)
func New() Service {
	return &memory{soldOut: map[string]bool{"sold-out": true}}
}

func (m *memory) Available(_ context.Context, item string) (bool, error) {
	return !m.soldOut[item], nil
}
