package http

import (
	"github.com/xgodev/boost/fx/modules/bootstrap/function/adapter/contrib/cloudevents"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = cloudevents.Root + ".http"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the nats adapter")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
