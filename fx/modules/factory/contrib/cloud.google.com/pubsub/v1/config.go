package pubsub

import (
	pubsubcfg "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"
)

// Root is the config root of the Pub/Sub factory fx modules.
const Root = pubsubcfg.Root

// IsEnabled returns true when the configured version is v1.
func IsEnabled() bool {
	return pubsubcfg.Enabled(pubsubcfg.VersionKey, pubsubcfg.V1)
}
