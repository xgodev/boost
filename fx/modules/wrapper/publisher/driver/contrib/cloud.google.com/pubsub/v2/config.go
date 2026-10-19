package pubsub

import (
	pubsubcfg "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"
	drivercfg "github.com/xgodev/boost/fx/modules/wrapper/publisher/driver/contrib/cloud.google.com/pubsub"
)

// Root is the config root of the Pub/Sub publisher driver fx modules.
const Root = drivercfg.Root

// IsEnabled returns true when the driver is enabled and its version is v2 (the default).
func IsEnabled() bool {
	return drivercfg.IsEnabled() && pubsubcfg.Enabled(drivercfg.VersionKey, pubsubcfg.V2)
}
