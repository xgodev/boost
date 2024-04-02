package logrus

import (
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus/v1"
)

// NewOptions returns options from config file or environment vars.
func NewOptions() (*logrus.Options, error) {
	return boost.NewOptionsWithPath[logrus.Options](root)
}
