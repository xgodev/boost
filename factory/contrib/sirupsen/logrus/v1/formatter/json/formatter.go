package json

import (
	"github.com/sirupsen/logrus"
	"github.com/xgodev/boost/wrapper/config"
)

// NewFormatter returns logrus formatter for json.
func NewFormatter() (logrus.Formatter, error) {
	return config.NewOptionsWithPath[logrus.JSONFormatter](root)
}
