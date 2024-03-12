package log

import (
	"github.com/xgodev/boost/factory/xgodev/boost.v1/grapper"
	"strings"

	"github.com/xgodev/boost/config"
)

const (
	root    = grapper.PluginsRoot + ".log"
	enabled = ".enabled"
)

func ConfigAdd(name string) {
	path := strings.Join([]string{root, ".", name}, "")
	config.Add(path+enabled, true, "enable/disable log grapper middleware")
}
