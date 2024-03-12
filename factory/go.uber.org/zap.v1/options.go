package zap

import (
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/log/contrib/go.uber.org/zap.v1"
)

// NewOptions returns configured zap logger options.
func NewOptions() (*zap.Options, error) {
	return factory.NewOptionsWithPath[zap.Options](root)
}
