// Package rest holds what the rest server of app api needs besides its
// handlers.
package rest

import (
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/multiserver/internal/app/api/background"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/health"
	multiserver "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/multi-server"
	status "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/model/restresponse"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/cors"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/gzip"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/requestid"
)

// Plugins are the usual plugins of a boost echo server plus the multi-server
// check route (GET /check), which answers 503 when one of the servers of the
// app stopped.
//
// It receives the background server only so the rest door builds it: a server
// that is not a door joins multiserver when something the app serves
// depends on it.
// @Provide(index=0)
func Plugins(_ *background.Server) []boostecho.Plugin {
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
