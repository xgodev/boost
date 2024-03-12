package log

import (
	"github.com/xgodev/boost/factory/local/wrapper"
	"strings"

	"github.com/xgodev/boost/config"
)

const (
	root    = wrapper.PluginsRoot + ".log"
	enabled = ".enabled"
)

func ConfigAdd(name string) {
	path := strings.Join([]string{root, ".", name}, "")
	config.Add(path+enabled, true, "enable/disable log grapper middleware")
}
