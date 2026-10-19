package boost

// Options configures the publisher of price changed events. inject generates
// its config keys (inject.examples.samename.publisher.pricechanged.boost.*)
// and provides *Options, apart from the Options of the other boost package.
// @Config
type Options struct {
	Subject string `default:"precos-alteracao" desc:"subject the price changed events go to"`
}
