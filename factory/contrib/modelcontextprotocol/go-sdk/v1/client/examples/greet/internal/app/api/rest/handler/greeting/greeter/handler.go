// Package greeter serves GET /greetings/:name.
package greeter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client/examples/greet/internal/pkg/domain/service/greeting/greeter"
)

// Response is the body of GET /greetings/:name.
type Response struct {
	Greeting string `json:"greeting"`
}

// Handler greets.
type Handler struct {
	svc greeter.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc greeter.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle answers with the greeting of the name.
// @REST(method=GET, path=/greetings/:name)
func (h *Handler) Handle(c echo.Context) error {
	greeting, err := h.svc.Greet(c.Request().Context(), c.Param("name"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{Greeting: greeting.Message})
}
