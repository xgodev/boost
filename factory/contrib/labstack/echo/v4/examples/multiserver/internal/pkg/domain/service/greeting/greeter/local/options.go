package local

// Options configures the greeting; inject generates its config key
// (internal.pkg.domain.service.greeting.greeter.local.message).
// @Config
type Options struct {
	Message string `default:"Hello World!!" desc:"message of the greeting"`
}
