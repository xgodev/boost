package datadog

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/aws/aws-sdk-go.v2"
)

const (
	root    = aws.PluginsRoot + ".datadog"
	enabled = ".enabled"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(enabled, true, "enable/disable datadog integration")
}
