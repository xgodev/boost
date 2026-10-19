package mongodb

// Options configures the mongo stock; inject generates its config keys
// (internal.pkg.datastore.stock.mongodb.*).
// @Config
type Options struct {
	Collection string `default:"stock" desc:"collection of the units in stock"`
}
