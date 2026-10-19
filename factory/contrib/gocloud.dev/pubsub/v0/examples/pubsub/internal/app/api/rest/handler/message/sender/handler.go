// Package sender serves POST /messages.
package sender

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0/examples/pubsub/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/gocloud.dev/pubsub/v0/examples/pubsub/internal/pkg/domain/service/message/sender"
	"github.com/xgodev/boost/model/errors"
)

// Body is the body of POST /messages.
type Body struct {
	Text string `json:"text"`
}

// Handler sends messages.
type Handler struct {
	svc sender.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc sender.Service) *Handler {
	return &Handler{svc: svc}
}

// Handle sends the message and answers 202.
// @REST(method=POST, path=/messages)
func (h *Handler) Handle(c echo.Context) error {
	var body Body
	if err := c.Bind(&body); err != nil {
		return errors.Wrap(err, errors.BadRequestf("invalid body"))
	}
	if err := h.svc.Send(c.Request().Context(), &entity.Message{Text: body.Text}); err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)
}
