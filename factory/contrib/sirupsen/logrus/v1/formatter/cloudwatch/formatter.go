package cloudwatch

import (
	"github.com/ravernkoh/cwlogsfmt"
	"github.com/sirupsen/logrus"
	"github.com/xgodev/boost/wrapper/config"
)

// NewFormatter returns logrus formatter for cloudwatch.
func NewFormatter() (logrus.Formatter, error) {
	return config.NewOptionsWithPath[cwlogsfmt.CloudWatchLogsFormatter](root)
}
