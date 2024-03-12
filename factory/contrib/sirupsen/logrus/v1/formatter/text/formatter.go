package text

import (
	"github.com/sirupsen/logrus"
	"github.com/xgodev/boost"
)

// NewFormatter returns logrus formatter for text.
func NewFormatter() (logrus.Formatter, error) {
	return boost.NewOptionsWithPath[logrus.TextFormatter](root)
}
