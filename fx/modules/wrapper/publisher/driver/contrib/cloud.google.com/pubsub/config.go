// Package pubsub declares, once, the config the Pub/Sub v1 and v2 publisher
// driver fx modules share: whether the driver is on and which version it is.
package pubsub

import (
	pubsubcfg "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"
	"github.com/xgodev/boost/fx/modules/wrapper/publisher/driver"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	// Root is the config root of the Pub/Sub publisher driver fx modules.
	Root    = driver.Root + ".pubsub"
	enabled = Root + ".enabled"
	// VersionKey chooses which driver fx module provides the publisher.Driver.
	VersionKey = Root + ".version"
)

func init() {
	config.Add(enabled, true, "enables/disables the pubsub driver")
	config.Add(VersionKey, pubsubcfg.V2, "defines the pubsub version: v1 or v2")
}

// IsEnabled returns true when the driver is enabled, whatever its version.
func IsEnabled() bool {
	return config.Bool(enabled)
}
