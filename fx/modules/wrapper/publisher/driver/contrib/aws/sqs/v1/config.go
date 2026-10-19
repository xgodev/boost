package sqs

import (
	"github.com/xgodev/boost/fx/modules/wrapper/publisher/driver"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = driver.Root + ".sqs"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the sqs driver")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
