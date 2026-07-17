package pubsub

import (
	"github.com/xgodev/boost/fx/modules/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = adapter.Root + ".pubsub"
	enabled = Root + ".enabled"
	version = Root + ".version"
)

func init() {
	config.Add(enabled, true, "enables/disables the gcp pubsub adapter")
	config.Add(version, "v1", "defines the pubsub version: v1 or v2")
}

// IsEnabled returns true when the adapter is enabled and the configured version is v2.
func IsEnabled() bool {
	return config.Bool(enabled) && config.String(version) == "v2"
}
