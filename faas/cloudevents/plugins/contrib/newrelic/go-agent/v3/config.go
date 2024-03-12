package newrelic

import (
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/faas/cloudevents"
)

const (
	root    = cloudevents.PluginsRoot + ".newrelic"
	enabled = root + ".enabled"
	txName  = root + ".txName"
)

func init() {
	config.Add(enabled, true, "enable/disable newrelic middleware")
	config.Add(txName, "changeme", "cloudevents newrelic middleware tx name")
}

// IsEnabled reports whether the NewRelic middleware is enabled in the configuration.
func IsEnabled() bool {
	return config.Bool(enabled)
}

// TxName returns the configured New Relic transaction name.
func TxName() string {
	return config.String(txName)
}
