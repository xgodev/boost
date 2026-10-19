package pubsub

import (
	adaptercfg "github.com/xgodev/boost/fx/modules/bootstrap/function/adapter/contrib/cloud.google.com/pubsub"
	pubsubcfg "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"
)

// Root is the config root of the Pub/Sub function adapter fx modules.
const Root = adaptercfg.Root

// IsEnabled returns true when the adapter is enabled and its version is v2 (the default).
func IsEnabled() bool {
	return adaptercfg.IsEnabled() && pubsubcfg.Enabled(adaptercfg.VersionKey, pubsubcfg.V2)
}
