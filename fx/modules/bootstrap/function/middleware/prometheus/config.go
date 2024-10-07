package ignore_errors

import (
	"github.com/xgodev/boost/bootstrap/function/middleware"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = middleware.Root + ".prometheus"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the publisher middleware")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
