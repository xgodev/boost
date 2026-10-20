// Package creator serves POST /stocks.
package creator

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/creator"
	"github.com/xgodev/boost/examples/stock/internal/pkg/server/rest/model"
	"github.com/xgodev/boost/model/errors"
)

// Handler puts a new item in stock.
type Handler struct {
	svc creator.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc creator.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle creates the stock in the body.
// @REST(method=POST, path=/stocks)
func (h *Handler) Handle(c echo.Context) error {
	var body model.Stock
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequest(err, "body is not a stock")
	}
	if err := h.svc.Create(c.Request().Context(), body.Entity()); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, body)
}
