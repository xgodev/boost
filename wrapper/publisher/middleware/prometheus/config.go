package prometheus

import (
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/publisher/middleware"
)

const (
	root        = middleware.Root + ".prometheus"
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
