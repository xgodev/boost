package goka

import (
	"github.com/xgodev/boost/wrapper/log"
)

type Logger struct {
	level string
}

// NewLogger returns a new logger.
func NewLogger() *Logger {
	return &Logger{level: LogLevel()}
}

func (s *Logger) Print(msgs ...interface{}) {
	s.log("%v", msgs)
}

func (s *Logger) Println(msgs ...interface{}) {
	s.log("%v", msgs)
}

func (s *Logger) Printf(msg string, args ...interface{}) {
	s.log(msg, args...)
}

func (s *Logger) Debugf(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}

func (s *Logger) log(format string, args ...interface{}) {
	switch s.level {
	case "INFO":
		log.Infof(format, args...)
	case "TRACE":
		log.Tracef(format, args...)
	default:
		log.Debugf(format, args...)
	}
}
