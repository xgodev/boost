package stripslashes

import (
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = chi.PluginsRoot + ".stripSlashes"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable stripSlashes middleware")
}
