// Package fetcher serves GET /google.
package fetcher

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/googlecall/internal/pkg/domain/service/page/fetcher"
)

// Response is the body of GET /google.
type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Handler fetches the page and tells how it answered.
type Handler struct {
	svc fetcher.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc fetcher.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle fetches the page and answers with its status.
// @REST(method=GET, path=/google)
func (h *Handler) Handle(c echo.Context) error {
	page, err := h.svc.Fetch(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{Message: "Hello Google!!", Status: page.Status})
}
