package checker

// Options configures the check; inject generates its config keys
// (internal.app.stockcheck.job.stock.checker.*) and provides *Options to New.
// @Config
type Options struct {
	SKUs []string `config:"skus" default:"tshirt,mug" desc:"skus the check looks at"`
}
