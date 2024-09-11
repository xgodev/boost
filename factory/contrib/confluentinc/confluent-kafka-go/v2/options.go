package confluent

import (
	"github.com/xgodev/boost/wrapper/config"
)

// Options kafka connection options.
type Options struct {
	Brokers string
	Log     struct {
		Level   string
		Enabled bool
	}
	Producer struct {
		Acks    int
		Timeout struct {
			Request int
			Message int
		}
		Batch struct {
			Size        int
			NumMessages int
		}
	}
	Consumer struct {
		Topics           []string
		GroupId          string
		AutoOffsetReset  string
		EnableAutoCommit bool
		Protocol         string
	}
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
