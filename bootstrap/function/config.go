package function

import (
	"github.com/xgodev/boost/bootstrap"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root = bootstrap.Root + ".function"
	def  = Root + ".default"
)

func init() {
	config.Add(def, "", "default cmd")
}

// Default returns the default cmd name from config.
func Default() string {
	return config.String(def)
}
