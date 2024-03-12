package elasticsearch

import "github.com/xgodev/boost/log"

type DebugLogger struct {
}

func (l *DebugLogger) Printf(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}
