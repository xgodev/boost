// Package incrementer serves POST /counters/:name.
package incrementer

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9/examples/cluster/internal/pkg/domain/service/counter/incrementer"
)

// Response is the body of POST /counters/:name.
type Response struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// Handler increments a counter.
type Handler struct {
	svc incrementer.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc incrementer.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle increments the counter and answers with its count.
// @REST(method=POST, path=/counters/:name)
func (h *Handler) Handle(c echo.Context) error {
	counter, err := h.svc.Increment(c.Request().Context(), c.Param("name"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{Name: counter.Name, Count: counter.Count})
}
