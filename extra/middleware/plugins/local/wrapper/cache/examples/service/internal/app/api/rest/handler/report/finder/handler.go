// Package finder serves GET /reports/:id.
package finder

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/cache/examples/service/internal/pkg/domain/service/report/finder"
	"github.com/xgodev/boost/model/errors"
)

// Response is the body of GET /reports/:id.
type Response struct {
	ID    string `json:"id"`
	Build int64  `json:"build"`
}

// Handler finds reports.
type Handler struct {
	svc finder.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc finder.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle answers the report of the id of the path, 400 for an id that is not
// a number.
// @REST(method=GET, path=/reports/:id)
func (h *Handler) Handle(c echo.Context) error {
	report, err := h.svc.Find(c.Request().Context(), c.Param("id"))
	if errors.IsNotValid(err) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{ID: report.ID, Build: report.Build})
}
