package pubsub

import (
	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root   = adapter.Root + ".pubsub"
	topics = root + ".topics"
)

func init() {
	config.Add(topics, []string{"changeme"}, "pubsub listener topics")
}
