// Package second holds what the second rest server of app api needs.
package second

import (
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/health"
	multiserver "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/multi-server"
	status "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/model/restresponse"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/cors"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/gzip"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/requestid"
)

// Plugins of server second (configured under api.server.second): the usual
// ones, and the multi-server route /check, which answers 503 when one of
// the servers of the app is down.
// @Provide(name=second, index=0)
func Plugins() []boostecho.Plugin {
	return []boostecho.Plugin{
		cors.Register,
		requestid.Register,
		gzip.Register,
		log.Register,
		status.Register,
		multiserver.Register,
		health.Register,
	}
}
