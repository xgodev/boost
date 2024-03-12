package log

import (
	"github.com/xgodev/boost/factory/go.uber.org/zap.v1"
	"github.com/xgodev/boost/factory/rs/zerolog.v1"
	"github.com/xgodev/boost/factory/sirupsen/logrus.v1"
	"github.com/xgodev/boost/log"
)

// New initializes the log according to the configured type and formatter.
func New() {
	switch Type() {
	case "NOOP":
		log.NewNoop()
	case "ZEROLOG":
		zerolog.NewLogger()
	case "ZAP":
		zap.NewLogger()
	default:
		logrus.NewLogger()
	}
}
