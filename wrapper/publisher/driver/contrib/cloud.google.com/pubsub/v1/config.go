package pubsub

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = "boost.wrapper.driver.pubsub"
	logRoot = ".log"
	level   = logRoot + ".level"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+level, "DEBUG", "defines log level")
}
