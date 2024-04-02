package zap

import (
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/log/contrib/go.uber.org/zap/v1"
)

// NewLogger return a new zap logger.
func NewLogger() log.Logger {
	options, err := NewOptions()
	if err != nil {
		panic(err)
	}
	return zap.NewLoggerWithOptions(options)
}
