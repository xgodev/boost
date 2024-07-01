package fx

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root     = "boost.factory.fx"
	logLevel = root + ".log.level"
)

func init() {
	config.Add(logLevel, "DEBUG", "define log level")
}

func LogLevel() string {
	return config.String(logLevel)
}
