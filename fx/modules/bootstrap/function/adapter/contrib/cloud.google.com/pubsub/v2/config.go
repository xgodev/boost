package pubsub

import (
	"github.com/xgodev/boost/fx/modules/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = adapter.Root + ".pubsub.v2"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, false, "enables/disables the gcp pubsub adapter")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
