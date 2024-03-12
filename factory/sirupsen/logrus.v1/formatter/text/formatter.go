package text

import (
	"github.com/sirupsen/logrus"
	"github.com/xgodev/boost/factory"
)

// NewFormatter returns logrus formatter for text.
func NewFormatter() (logrus.Formatter, error) {
	return factory.NewOptionsWithPath[logrus.TextFormatter](root)
}
