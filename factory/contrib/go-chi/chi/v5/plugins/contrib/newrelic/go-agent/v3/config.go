package newrelic

import (
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root               = chi.PluginsRoot + ".newrelic"
	enabled            = ".enabled"
	webResponseEnabled = ".webresponse.enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable newrelic middleware")
	config.Add(path+webResponseEnabled, true, "enable/disable newrelic web response")
}
