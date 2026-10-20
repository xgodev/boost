// Package reserver serves POST /stocks/:sku/reserve.
package reserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/reserver"
	"github.com/xgodev/boost/examples/stock/internal/pkg/server/rest/model"
)

// Handler reserves one unit of an item.
type Handler struct {
	svc reserver.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc reserver.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle reserves one unit of the sku in the path.
// @REST(method=POST, path=/stocks/:sku/reserve)
func (h *Handler) Handle(c echo.Context) error {
	left, err := h.svc.Reserve(c.Request().Context(), c.Param("sku"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.FromEntity(left))
}
