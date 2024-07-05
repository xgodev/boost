package fx

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root     = "boost.factory.fx"
	logLevel = Root + ".log.level"
)

func init() {
	config.Add(logLevel, "DEBUG", "define log level")
}

func LogLevel() string {
	return config.String(logLevel)
}
