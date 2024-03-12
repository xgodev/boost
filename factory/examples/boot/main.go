package main

import (
	"os"

	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/factory/go.uber.org/zap.v1"
)

func init() {
	os.Setenv("IGNITE_ZAP_CONSOLE_FORMATTER", "JSON")
	os.Setenv("IGNITE_ZAP_CONSOLE_LEVEL", "DEBUG")
}

func main() {
	factory.Boot()
	zap.NewLogger()
}
