package memory

// Options configures the memory stock. inject generates its config keys
// (builder.examples.simple.stock.memory.*) and provides *Options to New.
// @Config
type Options struct {
	SoldOut []string `default:"sold-out" desc:"items the memory stock is out of"`
}
