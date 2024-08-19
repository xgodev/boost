package publisher

import (
	"github.com/xgodev/boost/wrapper/config"
)

type Options struct {
	Subject    string
	Deadletter struct {
		Enabled bool
		Subject string
		Errors  []error
	}
	Retry struct {
		Enabled bool
		Backoff int
		Errors  []error
	}
}

func NewOptions() (*Options, error) {
	o := &Options{}
	err := config.UnmarshalWithPath(root, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
