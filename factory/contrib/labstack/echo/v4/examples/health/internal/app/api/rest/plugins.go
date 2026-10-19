// Package rest holds what the rest server of app api needs besides its
// handlers.
package rest

import (
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/health"
)

// Plugins serve the checkers added to extra/health on /health
// (boost.factory.echo.plugins.health.route).
// @Provide(index=0)
func Plugins() []boostecho.Plugin {
	return []boostecho.Plugin{health.Register}
}
