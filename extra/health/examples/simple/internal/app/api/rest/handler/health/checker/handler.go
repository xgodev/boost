// Package checker serves GET /health.
package checker

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/extra/health/examples/simple/internal/pkg/domain/service/health/checker"
)

// Check is one item of the body of GET /health.
type Check struct {
	Name  string `json:"name"`
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// Handler answers the health of the app.
type Handler struct {
	svc checker.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc checker.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle answers every check, with 503 when a required one failed.
// @REST(method=GET, path=/health)
func (h *Handler) Handle(c echo.Context) error {
	status := http.StatusOK
	var body []Check
	for _, check := range h.svc.Check(c.Request().Context()) {
		if check.Error != "" && check.Required {
			status = http.StatusServiceUnavailable
		}
		body = append(body, Check{Name: check.Name, OK: check.Error == "", Error: check.Error})
	}
	return c.JSON(status, body)
}
