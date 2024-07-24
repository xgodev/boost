package noop

import (
	"github.com/xgodev/boost/fx/modules/local/wrapper/publisher/driver"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = driver.Root + ".noop"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the noop driver")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
