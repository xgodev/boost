package azure

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options configures the Azure Functions adapter.
type Options struct {
	// Port is where the functions host calls the custom handler.
	Port int
	// Trigger is the binding name of the trigger in the invocation data;
	// empty takes the only input.
	Trigger string
	// Output is the output binding that receives the handler result;
	// $return is the function return value.
	Output string
	// RoutePrefix is the route prefix of the http functions, to tell a
	// forwarded http request from an invocation.
	RoutePrefix string
	// MaxWorkers bounds the items of a batch handled at the same time.
	MaxWorkers int
}

// NewOptions returns the options read from config.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}
