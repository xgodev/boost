// Package pubsub holds what the Pub/Sub v1 and v2 fx modules share: the
// version values and the factory's version key, declared once.
package pubsub

import (
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	// V1 selects cloud.google.com/go/pubsub.
	V1 = "v1"
	// V2 selects cloud.google.com/go/pubsub/v2, the default.
	V2 = "v2"

	// Root is the config root of the Pub/Sub factory fx modules.
	Root = "boost.factory.gcp.pubsub"
	// VersionKey chooses which factory fx module provides the client.
	VersionKey = Root + ".version"
)

func init() {
	config.Add(VersionKey, V2, "defines the pubsub version: v1 or v2")
}

// Version returns the Pub/Sub version configured at key: V1 or V2. Any other
// value is a NotValid error naming the key and the value.
func Version(key string) (string, error) {
	v := config.String(key)
	if v != V1 && v != V2 {
		return "", errors.NotValidf("%s is %q; valid values: %s, %s", key, v, V1, V2)
	}
	return v, nil
}

// Enabled reports whether key is configured as want; an invalid value is
// never enabled (Version reports it).
func Enabled(key, want string) bool {
	v, err := Version(key)
	return err == nil && v == want
}
