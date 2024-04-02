package hystrix

import (
	"strings"

	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/local/middleware"
)

const (
	root    = middleware.PluginsRoot + ".hystrix"
	enabled = ".enabled"
)

func ConfigAdd(name string) {
	path := strings.Join([]string{root, ".", name}, "")
	config.Add(path+enabled, true, "enable/disable hystrix grapper middleware")
}
