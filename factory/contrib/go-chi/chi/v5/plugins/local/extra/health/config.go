package health

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
)

const (
	root    = chi.PluginsRoot + ".health"
	enabled = ".enabled"
	route   = ".route"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable health route")
	config.Add(path+route, "/health", "define status url")
}
