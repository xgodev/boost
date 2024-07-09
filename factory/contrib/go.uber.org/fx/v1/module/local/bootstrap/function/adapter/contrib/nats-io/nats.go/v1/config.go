package nats

import (
	"github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/local/bootstrap/function/adapter"
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
