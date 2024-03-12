package godror

import (
	"github.com/xgodev/boost"
	"time"
)

// Options represents a godror options.
type Options struct {
	ConnectString       string
	Username            string
	Password            string
	MaxLifetime         time.Duration
	SessionTimeout      time.Duration
	WaitTimeout         time.Duration
	MaxSessions         int
	SessionIncrement    int
	MinSessions         int
	MaxSessionsPerShard int
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return boost.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return boost.NewOptionsWithPath[Options](root, path)
}
