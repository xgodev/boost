package log

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/faas/cloudevents"
)

const (
	root    = cloudevents.PluginsRoot + ".logger"
	enabled = root + ".enabled"
	level   = root + ".level"
)

func init() {
	config.Add(enabled, true, "enable/disable logger middleware")
	config.Add(level, "INFO", "sets log level INFO/DEBUG/TRACE")
}

// IsEnabled reports whether the logger middleware is enabled in the configuration.
func IsEnabled() bool {
	return config.Bool(enabled)
}

// Level returns the configured log level.
func Level() string {
	return config.String(level)
}
