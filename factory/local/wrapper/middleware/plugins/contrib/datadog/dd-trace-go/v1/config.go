package datadog

import (
	"github.com/xgodev/boost/wrapper/config"
	"strings"

	"github.com/xgodev/boost/factory/local/wrapper/middleware"
)

const (
	root    = middleware.PluginsRoot + ".datadog"
	enabled = ".enabled"
)

func ConfigAdd(name string) {
	path := strings.Join([]string{root, ".", name}, "")
	config.Add(path+enabled, true, "enable/disable datadog grapper middleware")
}
