package cache

import (
	"github.com/xgodev/boost/factory/local/wrapper/middleware"
	"github.com/xgodev/boost/wrapper/config"
	"strings"
)

const (
	root    = middleware.PluginsRoot + ".cache"
	enabled = ".enabled"
)

func ConfigAdd(name string) {
	path := strings.Join([]string{root, ".", name}, "")
	config.Add(path+enabled, true, "enable/disable cache grapper middleware")
}
