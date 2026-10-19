package boost

// Options configures the publisher of price received events. inject generates
// its config keys (inject.examples.samename.publisher.pricereceived.boost.*)
// and provides *Options, apart from the Options of the other boost package.
// @Config
type Options struct {
	Subject string `default:"precos-entrada" desc:"subject the price received events go to"`
}
