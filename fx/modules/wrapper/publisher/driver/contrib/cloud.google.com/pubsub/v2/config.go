package pubsub

import (
	"github.com/xgodev/boost/fx/modules/wrapper/publisher/driver"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = driver.Root + ".pubsub"
	enabled = Root + ".enabled"
	version = Root + ".version"
)

func init() {
	config.Add(enabled, true, "enables/disables the pubsub driver")
	config.Add(version, "v1", "defines the pubsub version: v1 or v2")
}

// IsEnabled returns true when the driver is enabled and the configured version is v2.
func IsEnabled() bool {
	return config.Bool(enabled) && config.String(version) == "v2"
}
