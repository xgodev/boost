package pubsub

import (
	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root          = adapter.Root + ".pubsub.v2"
	subscriptions = root + ".subscriptions"
	concurrency   = root + ".concurrency"
)

func init() {
	config.Add(subscriptions, []string{"changeme"}, "pubsub listener topics")
	config.Add(concurrency, 10, "pubsub max concurrent workers")
}
