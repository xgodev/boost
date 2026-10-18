// Package memory keeps the stock in memory: the default, no infra needed.
package memory

import (
	"context"
	"sync"

	"github.com/acme/shop/internal/pkg/domain/repository/stock"
	"github.com/xgodev/boost/model/errors"
)

type datastore struct {
	mu    sync.Mutex
	units map[string]int
}

// New provides the stock wired when no other is configured.
// @Provide(index=0)
// @Default
func New() stock.Repository {
	return &datastore{units: map[string]int{"tshirt": 3}}
}

func (d *datastore) Take(_ context.Context, sku string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.units[sku] == 0 {
		return errors.Conflictf("%s is sold out", sku)
	}
	d.units[sku]--
	return nil
}
