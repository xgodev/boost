package pubsub

import (
	"github.com/xgodev/boost/fx/modules/wrapper/publisher/driver"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = driver.Root + ".pubsub.v2"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the pubsub driver")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
