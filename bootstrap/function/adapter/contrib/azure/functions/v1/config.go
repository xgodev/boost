package azure

import (
	"os"
	"strconv"

	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = adapter.Root + ".azure"
	port        = root + ".port"
	trigger     = root + ".trigger"
	output      = root + ".output"
	routePrefix = root + ".routePrefix"
	maxWorkers  = root + ".maxWorkers"
)

func init() {
	config.Add(port, customHandlerPort(), "port the functions host calls; defaults to FUNCTIONS_CUSTOMHANDLER_PORT, set by the host")
	config.Add(trigger, "", "binding name of the trigger in the invocation data; empty takes the only input")
	config.Add(output, "$return", "output binding that receives the handler result; $return is the function return value")
	config.Add(routePrefix, "api", "route prefix of the http functions (host.json extensions.http.routePrefix), to tell a forwarded http request from an invocation")
	config.Add(maxWorkers, 10, "items of a batch invocation handled at the same time; 1 keeps their order")
}

// customHandlerPort is the port the functions host gives the custom handler
// in FUNCTIONS_CUSTOMHANDLER_PORT, or 8080 outside the host.
func customHandlerPort() int {
	if p, err := strconv.Atoi(os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")); err == nil && p > 0 {
		return p
	}
	return 8080
}
