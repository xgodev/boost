package nats

import (
	"github.com/xgodev/boost/fx/modules/local/wrapper/publisher/driver"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = driver.Root + ".nats"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the nats driver")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
