package boost

// Options configures the reservation event; inject generates its config keys
// (internal.pkg.publisher.stock.reserved.boost.*).
// @Config
type Options struct {
	Subject string `default:"stock.reserved" desc:"subject (topic) the reservation event is published to"`
}
