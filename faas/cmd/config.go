package cmd

import (
	"github.com/xgodev/boost/config"
)

const (
	root = "faas.cmd"
	def  = root + ".default"
)

func init() {
	config.Add(def, "", "default cmd")
}

// Default returns the default cmd name from config.
func Default() string {
	return config.String(def)
}
