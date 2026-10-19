package fx

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root     = "boost.factory.fx"
	logLevel = Root + ".log.level"
)

func init() {
	config.Add(logLevel, "DEBUG", "level of fx's lines (TRACE, DEBUG, INFO); an error that stops the app is always logged at ERROR")
}

func LogLevel() string {
	return config.String(logLevel)
}
