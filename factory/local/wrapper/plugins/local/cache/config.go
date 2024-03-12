package cache

import (
	"strings"

	"github.com/xgodev/boost/config"
)

const (
	root    = grapper.PluginsRoot + ".cache"
	enabled = ".enabled"
)

func ConfigAdd(name string) {
	path := strings.Join([]string{root, ".", name}, "")
	config.Add(path+enabled, true, "enable/disable cache grapper middleware")
}
