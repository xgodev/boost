package zerolog

import (
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/log/contrib/rs/zerolog/v1"
)

// NewOptions returns options from config file or environment vars.
func NewOptions() (*zerolog.zerolog, error) {
	return boost.NewOptionsWithPath[zerolog.Options](root)
}
