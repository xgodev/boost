package logrus

import (
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus/v1"
)

// NewOptions returns options from config file or environment vars.
func NewOptions() (*logrus.Options, error) {
	return config.NewOptionsWithPath[logrus.Options](root)
}
