package json

import (
	"github.com/sirupsen/logrus"
	"github.com/xgodev/boost"
)

// NewFormatter returns logrus formatter for json.
func NewFormatter() (logrus.Formatter, error) {
	return boost.NewOptionsWithPath[logrus.JSONFormatter](root)
}
