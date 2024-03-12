package cloudwatch

import (
	"github.com/ravernkoh/cwlogsfmt"
	"github.com/sirupsen/logrus"
	"github.com/xgodev/boost/factory"
)

// NewFormatter returns logrus formatter for cloudwatch.
func NewFormatter() (logrus.Formatter, error) {
	return factory.NewOptionsWithPath[cwlogsfmt.CloudWatchLogsFormatter](root)
}
