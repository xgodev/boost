package service

import (
	"github.com/dapr/go-sdk/service/common"
	"github.com/xgodev/boost/wrapper/config"
)

// Options can be used to create customized handler.
type Options struct {
	Subscriptions []common.Subscription
	Service       struct {
		Invocation struct {
			Enabled bool
			Name    string
		}
	}
}

// DefaultOptions returns options based in config.
func DefaultOptions() (*Options, error) {

	o := &Options{}

	err := config.UnmarshalWithPath(root, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}
