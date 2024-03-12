package hystrix

import (
	"strings"

	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/xgodev/boost.v1/grapper"
)

const (
	root    = grapper.PluginsRoot + ".hystrix"
	enabled = ".enabled"
)

func ConfigAdd(name string) {
	path := strings.Join([]string{root, ".", name}, "")
	config.Add(path+enabled, true, "enable/disable hystrix grapper middleware")
}
