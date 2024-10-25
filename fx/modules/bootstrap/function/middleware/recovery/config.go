package recovery

import (
	"github.com/xgodev/boost/bootstrap/function/middleware"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = middleware.Root + ".recovery"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the recovery middleware")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
