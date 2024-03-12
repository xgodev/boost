package json

import (
	"github.com/sirupsen/logrus"
	"github.com/xgodev/boost/factory"
)

// NewFormatter returns logrus formatter for json.
func NewFormatter() (logrus.Formatter, error) {
	return factory.NewOptionsWithPath[logrus.JSONFormatter](root)
}
