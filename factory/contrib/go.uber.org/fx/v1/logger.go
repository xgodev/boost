package fx

import (
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

// Logger represents a logger for fx.
type Logger struct {
	level string
}

// Printf logs format and args according to log level.
func (p *Logger) Printf(format string, args ...interface{}) {
	switch p.level {
	case "INFO":
		log.Infof(format, args...)
	case "TRACE":
		log.Tracef(format, args...)
	default:
		log.Debugf(format, args...)
	}
}

// NewLogger returns a new logger.
func NewLogger() fx.Printer {
	return &Logger{level: LogLevel()}
}
