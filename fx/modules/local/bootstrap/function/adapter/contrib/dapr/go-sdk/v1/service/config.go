package service

import (
	"github.com/xgodev/boost/fx/modules/local/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = adapter.Root + ".dapr.service"
	enabled = Root + ".enabled"
	tp      = Root + ".type"
)

func init() {
	config.Add(enabled, true, "enables/disables the dapr service adapter")
	config.Add(tp, "http", "sets http/grpc type of dapr service adapter")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}

func Type() string {
	return config.String(tp)
}
