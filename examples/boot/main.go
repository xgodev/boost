package main

import (
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.uber.org/zap/v1"
	"os"
)

func init() {
	os.Setenv("IGNITE_ZAP_CONSOLE_FORMATTER", "JSON")
	os.Setenv("IGNITE_ZAP_CONSOLE_LEVEL", "DEBUG")
}

func main() {
	boost.Start()
	zap.NewLogger()
}
