package greeting

// Options configures the default greeter. inject generates its config keys
// (inject.examples.simple.greeting.*) and provides *Options.
// @Config
type Options struct {
	Prefix string `default:"Hello" desc:"opening word of the default greeting"`
}
