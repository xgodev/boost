// Package catalog stores the prices of the catalog in Bigtable through the
// "catalog" client.
package catalog

import (
	"context"

	"cloud.google.com/go/bigtable"
)

const (
	// Table holds one row per SKU.
	Table = "prices"
	// Family is the column family of the price.
	Family = "p"
	column = "value"
)

// Prices reads and writes prices.
type Prices struct {
	table *bigtable.Table
}

// NewPrices provides the prices on the catalog client, which inject builds
// with the log plugin.
// @Provide(index=0)
// @Inject(index=0, name=catalog, plugins=[log])
func NewPrices(client *bigtable.Client) *Prices {
	return &Prices{table: client.Open(Table)}
}

// Set writes the price of sku.
func (p *Prices) Set(ctx context.Context, sku, price string) error {
	mut := bigtable.NewMutation()
	mut.Set(Family, column, bigtable.Now(), []byte(price))
	return p.table.Apply(ctx, sku, mut)
}

// Get reads the price of sku, empty when it has none.
func (p *Prices) Get(ctx context.Context, sku string) (string, error) {
	row, err := p.table.ReadRow(ctx, sku, bigtable.RowFilter(bigtable.LatestNFilter(1)))
	if err != nil {
		return "", err
	}
	for _, item := range row[Family] {
		return string(item.Value), nil
	}
	return "", nil
}
