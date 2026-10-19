// Package describer serves GET /application.
package describer

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config/examples/env_file/internal/pkg/domain/service/application/describer"
)

// Handler describes the application.
type Handler struct {
	svc describer.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc describer.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle answers the settings of the application.
// @REST(method=GET, path=/application)
func (h *Handler) Handle(c echo.Context) error {
	app, err := h.svc.Describe(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response(app))
}
