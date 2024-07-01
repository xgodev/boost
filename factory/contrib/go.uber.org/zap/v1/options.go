package zap

import (
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log/contrib/go.uber.org/zap/v1"
)

// NewOptions returns configured zap logger options.
func NewOptions() (*zap.Options, error) {
	return config.NewOptionsWithPath[zap.Options](root)
}
