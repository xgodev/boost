package function

import (
	"github.com/xgodev/boost/bootstrap"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	Root = bootstrap.Root + ".function"
	def  = Root + ".default"
	name = Root + ".name"
)

func init() {
	config.Add(def, "", "default cmd")
	config.Add(name, "func", "func name")
}

// Default returns the default cmd name from config.
func Default() string {
	return config.String(def)
}

// Name returns func name
func Name() string {
	return config.String(name)
}
