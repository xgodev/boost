package prometheus

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = "boost.factory.prometheus"
	pushGateway = root + ".pushGateway"
	pgEnabled   = pushGateway + ".enabled"
	pgURL       = pushGateway + ".url"
	pgInterval  = pushGateway + ".interval" // new
)

func init() {
	config.Add(pgEnabled, false, "enable/disable prometheus push gateway")
	config.Add(pgURL, "http://localhost:9091", "prometheus push gateway URL")
	config.Add(pgInterval, "10s", "interval between push gateway pushes") // default 10 seconds
}

var (
	pgEnabledVar  *bool
	pgURLVar      *string
	pgIntervalVar *time.Duration
)

func PushGatewayEnabled() bool {
	if pgEnabledVar == nil {
		v := config.Bool(pgEnabled)
		pgEnabledVar = &v
	}
	return *pgEnabledVar
}

func PushGatewayURL() string {
	if pgURLVar == nil {
		v := config.String(pgURL)
		pgURLVar = &v
	}
	return *pgURLVar
}

// PushInterval returns the duration between pushes
func PushInterval() time.Duration {
	if pgIntervalVar == nil {
		d := config.Duration(pgInterval)
		pgIntervalVar = &d
	}
	return *pgIntervalVar
}
