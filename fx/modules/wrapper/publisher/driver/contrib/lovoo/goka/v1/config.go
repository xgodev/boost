package goka

import (
	"github.com/xgodev/boost/fx/modules/wrapper/publisher/driver"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = driver.Root + ".kafka_goka"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the goka driver")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
