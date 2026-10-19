// Package pinger serves GET /ping.
package pinger

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1/examples/ping/internal/pkg/domain/service/database/pinger"
)

// Response is the body of GET /ping.
type Response struct {
	Database string `json:"database"`
	Status   string `json:"status"`
}

// Handler pings the database.
type Handler struct {
	svc pinger.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc pinger.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle pings the database and answers with its status.
// @REST(method=GET, path=/ping)
func (h *Handler) Handle(c echo.Context) error {
	db, err := h.svc.Ping(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{Database: db.Name, Status: db.Status})
}
