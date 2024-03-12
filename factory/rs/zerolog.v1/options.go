package zerolog

import (
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/log/contrib/rs/zerolog.v1"
)

// NewOptions returns options from config file or environment vars.
func NewOptions() (*zerolog.Options, error) {
	return factory.NewOptionsWithPath[zerolog.Options](root)
}
