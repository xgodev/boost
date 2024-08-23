package confluent

import (
	"github.com/xgodev/boost/fx/modules/wrapper/publisher/driver"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = driver.Root + ".kafka_confluent"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the confluent driver")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
