package otel

import (
	"github.com/go-logr/logr"
	"github.com/xgodev/boost/wrapper/log"
)

type Logger struct {
}

func (l Logger) Init(info logr.RuntimeInfo) {
}

func (l Logger) Enabled(level int) bool {
	return true
}

func (l Logger) Info(level int, msg string, keysAndValues ...any) {
	log.Infof(msg, keysAndValues...)
}

func (l Logger) Error(err error, msg string, keysAndValues ...any) {
	log.Errorf(msg, keysAndValues...)
}

func (l Logger) WithValues(keysAndValues ...any) logr.LogSink {
	return l
}

func (l Logger) WithName(name string) logr.LogSink {
	return l
}
