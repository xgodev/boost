package logrus

import (
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/log/contrib/sirupsen/logrus.v1"
)

// NewOptions returns options from config file or environment vars.
func NewOptions() (*logrus.Options, error) {
	return factory.NewOptionsWithPath[logrus.Options](root)
}
