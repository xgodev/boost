package kafka

import (
	"github.com/xgodev/boost/fx/modules/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = adapter.Root + ".kafkago"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the kafka-go adapter")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
