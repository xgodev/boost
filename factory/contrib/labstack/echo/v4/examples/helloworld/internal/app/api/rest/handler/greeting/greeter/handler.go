// Package greeter serves GET /hello-world.
package greeter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/helloworld/internal/pkg/domain/service/greeting/greeter"
)

// Response is the body of GET /hello-world.
type Response struct {
	Message string `json:"message"`
}

// Handler answers with the greeting of the service.
type Handler struct {
	svc greeter.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc greeter.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle answers with the greeting.
// @REST(method=GET, path=/hello-world)
func (h *Handler) Handle(c echo.Context) error {
	greeting, err := h.svc.Greet(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{Message: greeting.Message})
}
