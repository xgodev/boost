package pubsub

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root          = "boost.factory.pubsub"
	client        = Root + ".client"
	clientVersion = client + ".version"
)

func init() {
	config.Add(clientVersion, "v1", "changes the pubsub client version | v1 or v2")
}

func ClientVersion() string {
	return config.String(clientVersion)
}
