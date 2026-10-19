package elasticsearch

import "github.com/xgodev/boost/wrapper/log"

// DebugLogger sends the bulk indexer's debug lines to the boost logger.
type DebugLogger struct {
}

// Printf logs one debug line.
func (l *DebugLogger) Printf(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}
