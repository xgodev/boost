package oracle

// Options configures the oracle stock; inject generates its config keys
// (internal.pkg.datastore.stock.oracle.*).
// @Config
type Options struct {
	Table string `default:"stock" desc:"table of the units in stock: sku VARCHAR2 primary key, units NUMBER"`
}
