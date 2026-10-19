package stripslashes

import (
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = chi.PluginsRoot + ".stripslashes"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the strip slashes middleware plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable strip slashes middleware")
}
