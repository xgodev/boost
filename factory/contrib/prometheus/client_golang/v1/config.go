package prometheus

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = "boost.factory.prometheus"
	pushGateway = root + ".pushGateway"
	pgEnabled   = pushGateway + ".enabled"
	pgURL       = pushGateway + ".url"
	pgAsync     = pushGateway + ".async"
)

func init() {
	config.Add(pgEnabled, false, "enables/disables prometheus push gateway")
	config.Add(pgURL, "http://localhost", "defines prometheus push gateway url")
	config.Add(pgAsync, true, "enables/disables prometheus push gateway async")
}

var pgEnabledVar *bool
var pgAsyncVar *bool
var pgURLVar *string

func PushGatewayEnabled() bool {
	if pgEnabledVar == nil {
		chk := config.Bool(pgEnabled)
		pgEnabledVar = &chk
	}
	return *pgEnabledVar
}

func PushGatewayAsync() bool {
	if pgAsyncVar == nil {
		chk := config.Bool(pgAsync)
		pgAsyncVar = &chk
	}
	return *pgAsyncVar
}

func PushGatewayURL() string {
	if pgURLVar == nil {
		chk := config.String(pgURL)
		pgURLVar = &chk
	}
	return *pgURLVar
}
