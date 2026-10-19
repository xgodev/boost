package log

import (
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = chi.PluginsRoot + ".logger"
	enabled = ".enabled"
	level   = ".level"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the request log plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable logging request middleware")
	config.Add(path+level, "DEBUG", "sets log level INFO/DEBUG/TRACE")
}
