// Package rest holds what the rest server of app api needs besides its
// handlers.
package rest

import (
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	prometheus "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/contrib/prometheus/client_golang/v1"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/extra/error_handler"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/cors"
)

// Plugins add CORS, request logs and prometheus metrics, and turn boost
// errors into HTTP statuses: a NotFound answers 404.
// @Provide(index=0)
func Plugins() []boostecho.Plugin {
	return []boostecho.Plugin{
		cors.Register,
		log.Register,
		prometheus.Register,
		error_handler.Register,
	}
}
