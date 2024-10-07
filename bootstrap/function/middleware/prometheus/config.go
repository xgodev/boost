package prometheus

import (
	"github.com/xgodev/boost/bootstrap/function/middleware"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root         = middleware.Root + ".prometheus"
	functionName = root + ".functionName"
	pushGateway  = root + ".pushGateway"
	pgEnabled    = pushGateway + ".enabled"
	pgURL        = pushGateway + ".url"
	pgAsync      = pushGateway + ".async"
)

func init() {
	config.Add(functionName, "changeme", "defines prometheus function name")
	config.Add(pgEnabled, false, "enables/disables prometheus push gateway")
	config.Add(pgURL, "http://localhost", "defines prometheus push gateway url")
	config.Add(pgAsync, true, "enables/disables prometheus push gateway async")
}
