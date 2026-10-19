package multiserver

import (
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = chi.PluginsRoot + ".multiserver"
	enabled = ".enabled"
	route   = ".route"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the multi server check route plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable multi server check route")
	config.Add(path+route, "/check", "define multi server check url")
}
