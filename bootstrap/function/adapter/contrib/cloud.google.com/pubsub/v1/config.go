package pubsub

import (
	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

const (
	root          = adapter.Root + ".pubsub"
	subscriptions = root + ".subscriptions"
	backoff       = root + ".backoff"
	backoffBase   = root + ".backoffBase"
	maxBackoff    = root + ".maxBackoff"
	concurrency   = root + ".concurrency"
	retryLimit    = root + ".retryLimit"
)

func init() {
	config.Add(subscriptions, []string{"changeme"}, "pubsub listener topics")
	config.Add(backoff, true, "pubsub backoff")
	config.Add(backoffBase, 1*time.Second, "pubsub backoff base")
	config.Add(maxBackoff, 20*time.Second, "pubsub max backoff")
	config.Add(retryLimit, -1, "pubsub retry limit")
	config.Add(concurrency, 10, "pubsub retry limit")
}
