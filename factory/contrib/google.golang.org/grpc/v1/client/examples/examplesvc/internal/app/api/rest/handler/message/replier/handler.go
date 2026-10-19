// Package replier serves POST /messages.
package replier

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client/examples/examplesvc/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client/examples/examplesvc/internal/pkg/domain/service/message/replier"
	"github.com/xgodev/boost/model/errors"
)

// Body is the body of POST /messages and of its answer.
type Body struct {
	Message string `json:"message"`
}

// Handler sends messages.
type Handler struct {
	svc replier.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc replier.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle sends the message and answers with the reply.
// @REST(method=POST, path=/messages)
func (h *Handler) Handle(c echo.Context) error {
	var body Body
	if err := c.Bind(&body); err != nil {
		return errors.Wrap(err, errors.BadRequestf("invalid body"))
	}
	reply, err := h.svc.Reply(c.Request().Context(), &entity.Message{Text: body.Message})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Body{Message: reply.Text})
}
