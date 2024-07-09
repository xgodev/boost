package pubsub

import (
	"github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/local/bootstrap/function/middleware/publisher/driver"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = driver.Root + ".pubsub"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the nats driver")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
