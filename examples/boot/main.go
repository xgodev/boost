package main

import (
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.uber.org/zap/v1"
	"github.com/xgodev/boost/wrapper/log"
	"os"
)

func init() {
	os.Setenv("BOOST_FACTORY_ZAP_CONSOLE_FORMATTER", "JSON")
	os.Setenv("BOOST_FACTORY_ZAP_CONSOLE_LEVEL", "DEBUG")
}

func main() {
	boost.Start()
	zap.NewLogger()

	log.Infof("hello world!!")
}
