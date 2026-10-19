// Package pubsub declares, once, the config the Pub/Sub v1 and v2 function
// adapter fx modules share: whether the adapter is on and which version it is.
package pubsub

import (
	"github.com/xgodev/boost/fx/modules/bootstrap/function/adapter"
	pubsubcfg "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	// Root is the config root of the Pub/Sub function adapter fx modules.
	Root    = adapter.Root + ".pubsub"
	enabled = Root + ".enabled"
	// VersionKey chooses which adapter fx module subscribes.
	VersionKey = Root + ".version"
)

func init() {
	config.Add(enabled, true, "enables/disables the gcp pubsub adapter")
	config.Add(VersionKey, pubsubcfg.V2, "defines the pubsub version: v1 or v2")
}

// IsEnabled returns true when the adapter is enabled, whatever its version.
func IsEnabled() bool {
	return config.Bool(enabled)
}
