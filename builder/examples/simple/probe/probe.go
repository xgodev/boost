// Package probe answers the liveness probe on the admin server of app api,
// apart from the public routes.
package probe

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Live answers the liveness probe. Each @REST line is a route.
// @REST(method=GET, path=/live, server=admin)
// @REST(method=HEAD, path=/live, server=admin)
func Live(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
