package noop

import (
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/publisher/middleware"
)

const (
	Root    = middleware.Root + ".prometheus"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the prometheus middleware")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
