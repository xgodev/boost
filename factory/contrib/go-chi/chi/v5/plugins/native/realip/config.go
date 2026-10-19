package realip

import (
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = chi.PluginsRoot + ".realip"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the real ip middleware plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable real ip middleware; trusts X-Real-IP and X-Forwarded-For, so only behind a proxy that sets them")
}
