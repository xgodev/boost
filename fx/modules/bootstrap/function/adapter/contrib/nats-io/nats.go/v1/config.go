package nats

import (
	"github.com/xgodev/boost/fx/modules/local/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = adapter.Root + ".nats"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the nats adapter")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
