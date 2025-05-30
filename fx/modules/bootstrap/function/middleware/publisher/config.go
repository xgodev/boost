package publisher

import (
	"github.com/xgodev/boost/bootstrap/function/middleware"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = middleware.Root + ".publisher"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the publisher middleware")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
