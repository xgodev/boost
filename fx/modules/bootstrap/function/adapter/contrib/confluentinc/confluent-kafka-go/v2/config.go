package confluent

import (
	"github.com/xgodev/boost/fx/modules/local/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = adapter.Root + ".kafka_confluent"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the kafka confluent adapter")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
