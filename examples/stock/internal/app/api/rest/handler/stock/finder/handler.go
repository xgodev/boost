// Package finder serves GET /stocks/:sku.
package finder

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/finder"
	"github.com/xgodev/boost/examples/stock/internal/pkg/server/rest/model"
)

// Handler reads the stock of an item.
type Handler struct {
	svc finder.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc finder.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle returns the stock of the sku in the path.
// @REST(method=GET, path=/stocks/:sku)
func (h *Handler) Handle(c echo.Context) error {
	s, err := h.svc.Find(c.Request().Context(), c.Param("sku"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.FromEntity(s))
}
