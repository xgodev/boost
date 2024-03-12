package datadog

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
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
