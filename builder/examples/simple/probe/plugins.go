package probe

import (
	"context"

	"github.com/labstack/echo/v4"
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
)

// Plugins builds the admin server: every response says which server served
// it. A []Plugin provided under the name of a server reaches only that
// server; rest and grpc take the default (unnamed) one.
// @Provide(name=admin, index=0)
func Plugins() []boostecho.Plugin {
	return []boostecho.Plugin{servedBy("admin")}
}

func servedBy(server string) boostecho.Plugin {
	return func(_ context.Context, srv *boostecho.Server) error {
		srv.Instance().Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set("X-Served-By", server)
				return next(c)
			}
		})
		return nil
	}
}
