package log

import (
	"github.com/xgodev/boost/factory/contrib/go.uber.org/zap/v1"
	"github.com/xgodev/boost/factory/contrib/rs/zerolog/v1"
	"github.com/xgodev/boost/factory/contrib/sirupsen/logrus/v1"
	"github.com/xgodev/boost/wrapper/log"
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
