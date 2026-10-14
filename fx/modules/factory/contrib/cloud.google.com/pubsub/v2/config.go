package pubsub

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = "boost.factory.gcp.pubsub"
	version = Root + ".version"
)

func init() {
	config.Add(version, "v1", "defines the pubsub version: v1 or v2")
}

// IsEnabled returns true when the configured version is v2.
func IsEnabled() bool {
	return config.String(version) == "v2"
}
