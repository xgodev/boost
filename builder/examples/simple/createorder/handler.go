// Package createorder serves the order routes of app api.
package createorder

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/builder/examples/simple/order"
)

// Handler serves the order routes.
type Handler struct {
	svc *order.Service
}

// New provides the handler; the generated rest door gets it from here.
// @Provide(index=0)
func New(svc *order.Service) *Handler {
	return &Handler{svc: svc}
}

type request struct {
	Item string `json:"item"`
}

// Create stores the order in the body.
// @REST(method=POST, path=/orders)
func (h *Handler) Create(c echo.Context) error {
	var req request
	if err := c.Bind(&req); err != nil || req.Item == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "item is required")
	}
	return c.JSON(http.StatusCreated, h.svc.Create(req.Item))
}

// Get returns an order.
// @REST(method=GET, path=/orders/:id)
func (h *Handler) Get(c echo.Context) error {
	o, ok := h.svc.Get(c.Param("id"))
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}
	return c.JSON(http.StatusOK, o)
}
