package logger

import (
	"github.com/xgodev/boost/bootstrap/function/middleware"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root       = middleware.Root + ".logger"
	level      = root + ".level"
	errorStack = root + ".errorStack"
)

func init() {
	config.Add(level, "INFO", "defines log level")
	config.Add(errorStack, false, "defines if error stack should be logged")
}
