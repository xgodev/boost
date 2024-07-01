package zerolog

import (
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log/contrib/rs/zerolog/v1"
)

// NewOptions returns options from config file or environment vars.
func NewOptions() (*zerolog.Options, error) {
	return config.NewOptionsWithPath[zerolog.Options](root)
}
