package publisher

import (
	"github.com/xgodev/boost/wrapper/config"
)

type DeadletterOptions struct {
	Enabled bool
	Subject string
	Errors  []string
}

type Options struct {
	Subject    string
	Deadletter DeadletterOptions
}

func NewOptions() (*Options, error) {
	o := &Options{}
	err := config.UnmarshalWithPath(root, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
