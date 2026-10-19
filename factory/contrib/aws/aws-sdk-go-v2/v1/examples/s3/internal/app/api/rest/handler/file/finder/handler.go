// Package finder serves GET /files/:name.
package finder

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/examples/s3/internal/pkg/domain/service/file/finder"
)

// Response is the body of GET /files/:name.
type Response struct {
	Name         string    `json:"name"`
	LastModified time.Time `json:"lastModified"`
}

// Handler finds files.
type Handler struct {
	svc finder.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc finder.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle finds the file and answers when it was last modified.
// @REST(method=GET, path=/files/:name)
func (h *Handler) Handle(c echo.Context) error {
	file, err := h.svc.Find(c.Request().Context(), c.Param("name"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{Name: file.Name, LastModified: file.LastModified})
}
