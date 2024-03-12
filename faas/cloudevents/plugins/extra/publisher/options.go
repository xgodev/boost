package publisher

import (
	"github.com/xgodev/boost"
)

type Options struct {
	Enabled bool
	Success struct {
		Enabled bool
	}
	Error struct {
		Enabled bool
		Topic   string
	}
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return boost.NewOptionsWithPath[Options](root)
}
