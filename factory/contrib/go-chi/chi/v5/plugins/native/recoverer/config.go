package recoverer

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
)

const (
	root    = chi.PluginsRoot + ".recoverer"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable recoverer middleware")
}
