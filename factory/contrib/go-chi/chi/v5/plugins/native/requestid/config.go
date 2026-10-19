package requestid

import (
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = chi.PluginsRoot + ".requestid"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the request id middleware plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable request id middleware")
}
