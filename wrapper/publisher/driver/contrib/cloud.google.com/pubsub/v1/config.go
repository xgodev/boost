package pubsub

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = "boost.wrapper.publisher.driver.pubsub"
	logRoot     = ".log"
	orderingKey = ".orderingKey"
	level       = logRoot + ".level"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+level, "DEBUG", "defines log level")
	config.Add(path+orderingKey, false, "defines ordering key")
}
