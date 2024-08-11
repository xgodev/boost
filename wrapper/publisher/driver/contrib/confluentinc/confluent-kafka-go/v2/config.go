package confluent

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = "boost.wrapper.driver.kafka_confluent"
	logRoot = ".log"
	enabled = logRoot + ".enabled"
	level   = logRoot + ".level"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enables logging")
	config.Add(path+level, "DEBUG", "defines log level")
}
