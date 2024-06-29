package bigcache

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

type Options struct {
	Shards             int
	LifeWindow         time.Duration
	CleanWindow        time.Duration
	MaxEntriesInWindow int
	MaxEntrySize       int
	Verbose            bool
	HardMaxCacheSize   int
	StatsEnabled       bool
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
