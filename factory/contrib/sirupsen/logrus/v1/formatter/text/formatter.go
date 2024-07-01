package text

import (
	"github.com/sirupsen/logrus"
	"github.com/xgodev/boost/wrapper/config"
)

// NewFormatter returns logrus formatter for text.
func NewFormatter() (logrus.Formatter, error) {
	return config.NewOptionsWithPath[logrus.TextFormatter](root)
}
