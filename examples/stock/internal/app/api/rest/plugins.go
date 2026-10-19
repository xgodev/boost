// Package rest holds what the rest server of app api needs besides its
// handlers.
package rest

import (
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	prometheus "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/contrib/prometheus/client_golang/v1"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/extra/error_handler"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/health"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/recover"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/requestid"
)

// Plugins build the rest server: recover from panics, a request id per
// call, a log line per call, boost errors as HTTP statuses (NotFound 404,
// Conflict 409), GET /health with the checkers the clients register and
// GET /metrics. Each has its own enabled key under boost.factory.echo.plugins.
// @Provide(index=0)
func Plugins() []boostecho.Plugin {
	return []boostecho.Plugin{
		recover.Register,
		requestid.Register,
		log.Register,
		error_handler.Register,
		health.Register,
		prometheus.Register,
	}
}
