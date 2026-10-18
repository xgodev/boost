// Package reservestock serves POST /stock/:sku/reserve.
package reservestock

import (
	"net/http"

	"github.com/acme/shop/internal/pkg/domain/service/reservestock"
	"github.com/labstack/echo/v4"
)

// Handler reserves one unit of an item.
type Handler struct {
	svc reservestock.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc reservestock.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle reserves one unit of the sku in the path.
// @REST(method=POST, path=/stock/:sku/reserve)
func (h *Handler) Handle(c echo.Context) error {
	if err := h.svc.Reserve(c.Request().Context(), c.Param("sku")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
