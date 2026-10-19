package lambda

import (
	"github.com/xgodev/boost/fx/modules/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = adapter.Root + ".lambda"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the aws lambda adapter")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
