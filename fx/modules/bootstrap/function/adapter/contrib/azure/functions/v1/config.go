package azure

import (
	"github.com/xgodev/boost/fx/modules/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root    = adapter.Root + ".azure"
	enabled = Root + ".enabled"
)

func init() {
	config.Add(enabled, true, "enables/disables the azure functions adapter")
}

func IsEnabled() bool {
	return config.Bool(enabled)
}
