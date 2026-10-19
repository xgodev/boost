// Package finder serves GET /items/:id.
package finder

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/notfound/internal/pkg/domain/service/item/finder"
)

// Handler answers with the name of an item.
type Handler struct {
	svc finder.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc finder.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle answers with the name of the item in the path. The error of an
// unknown item is a NotFound, which the error_handler plugin answers 404.
// @REST(method=GET, path=/items/:id)
func (h *Handler) Handle(c echo.Context) error {
	item, err := h.svc.Find(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, item.Name)
}
