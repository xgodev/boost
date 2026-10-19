package memory

// Options configures the memory stock; inject generates its config keys
// (internal.pkg.datastore.stock.memory.*) and provides *Options to New.
// @Config
type Options struct {
	Seed []string `default:"tshirt=3,mug=1" desc:"sku=units the memory stock starts with"`
}
