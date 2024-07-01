package publisher

import (
	"github.com/xgodev/boost/bootstrap/cloudevents"
	"github.com/xgodev/boost/config"
)

const (
	root           = cloudevents.PluginsRoot + ".publisher"
	enabled        = root + ".enabled"
	successEnabled = root + ".success.enabled"
	errorEnabled   = root + ".error.enabled"
	errorTopic     = root + ".error.topic"
)

func init() {
	config.Add(enabled, true, "enable/disable publisher middleware")
	config.Add(successEnabled, false, "enable/disable success publisher middleware")
	config.Add(errorEnabled, false, "enable/disable error publisher middleware")
	config.Add(errorTopic, "changeme", "sets error topic publisher middleware")
}
