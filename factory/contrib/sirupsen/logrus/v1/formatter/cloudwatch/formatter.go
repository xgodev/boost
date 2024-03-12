package cloudwatch

import (
	"github.com/ravernkoh/cwlogsfmt"
	"github.com/sirupsen/logrus"
	"github.com/xgodev/boost"
)

// NewFormatter returns logrus formatter for cloudwatch.
func NewFormatter() (logrus.Formatter, error) {
	return boost.NewOptionsWithPath[cwlogsfmt.CloudWatchLogsFormatter](root)
}
