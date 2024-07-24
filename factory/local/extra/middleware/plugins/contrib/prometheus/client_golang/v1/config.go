package prometheus

import (
	"github.com/xgodev/boost/wrapper/config"
	"strings"

	"github.com/xgodev/boost/factory/local/extra/middleware"
)

const (
	root    = middleware.PluginsRoot + ".prometheus"
	enabled = ".enabled"
)

func ConfigAdd(name string) {
	path := strings.Join([]string{root, ".", name}, "")
	config.Add(path+enabled, true, "enable/disable prometheus grapper middleware")
}
