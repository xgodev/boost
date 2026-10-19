package kafka

import (
	"github.com/xgodev/boost/fx/modules/wrapper/publisher/driver"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = driver.Root + ".kafkago"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the kafka-go driver")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
