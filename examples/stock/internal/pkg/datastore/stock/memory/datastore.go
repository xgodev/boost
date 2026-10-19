// Package memory keeps the stock in memory: the default, no infra needed.
package memory

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/repository/stock"
	"github.com/xgodev/boost/model/errors"
)

type datastore struct {
	mu    sync.Mutex
	units map[string]int64
}

// New provides the stock wired when no other is configured, holding the
// items of Options.Seed.
// @Provide(index=0)
// @Default
func New(o *Options) (stock.Repository, error) {
	units := map[string]int64{}
	for _, entry := range o.Seed {
		sku, n, ok := strings.Cut(entry, "=")
		count, err := strconv.ParseInt(n, 10, 64)
		if !ok || sku == "" || err != nil || count < 0 {
			return nil, errors.NotValidf("seed %q is not sku=units", entry)
		}
		units[sku] = count
	}
	return &datastore{units: units}, nil
}

func (d *datastore) Create(_ context.Context, s *entity.Stock) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.units[s.SKU]; ok {
		return errors.AlreadyExistsf("%s is already in stock", s.SKU)
	}
	d.units[s.SKU] = s.Units
	return nil
}

func (d *datastore) Get(_ context.Context, sku string) (*entity.Stock, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	units, ok := d.units[sku]
	if !ok {
		return nil, errors.NotFoundf("%s is not in stock", sku)
	}
	return &entity.Stock{SKU: sku, Units: units}, nil
}

func (d *datastore) Take(_ context.Context, sku string) (*entity.Stock, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	units, ok := d.units[sku]
	if !ok {
		return nil, errors.NotFoundf("%s is not in stock", sku)
	}
	if units == 0 {
		return nil, errors.Conflictf("%s is sold out", sku)
	}
	d.units[sku] = units - 1
	return &entity.Stock{SKU: sku, Units: units - 1}, nil
}
