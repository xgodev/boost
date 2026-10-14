package datadog

import (
	ddtrace "github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

// Logger represents Datadog's Logger implementation.
type Logger struct {
}

// NewLogger returns implementation of Datadog's Logger interface.
func NewLogger() ddtrace.Logger {
	return &Logger{}
}

// Log  logs msg according to logLevel
func (l *Logger) Log(msg string) {

	var fn func(args ...interface{})

	switch config.String(logLevel) {
	case "INFO":
		fn = log.Info
	case "DEBUG":
		fn = log.Debug
	default:
		fn = log.Debug
	}

	fn(msg)
}
