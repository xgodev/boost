package datadog

import (
	girest "github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root          = girest.PluginsRoot + ".datadog"
	operationName = ".operationName"
	enabled       = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+operationName, "http.request", "defines span operation name")
	config.Add(path+enabled, true, "enable/disable datadog integration")
}
