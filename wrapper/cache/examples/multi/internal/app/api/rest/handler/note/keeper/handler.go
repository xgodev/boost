// Package keeper serves the notes: PUT and GET /notes/:key.
package keeper

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/cache/examples/multi/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/wrapper/cache/examples/multi/internal/pkg/domain/service/note/keeper"
)

// Request is the body of PUT /notes/:key.
type Request struct {
	Text string `json:"text"`
}

// Response is the body of GET /notes/:key.
type Response struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

// Handler keeps notes.
type Handler struct {
	svc keeper.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc keeper.Service) *Handler {
	return &Handler{svc: svc}
}

// Save keeps the note of the body under the key of the path.
// @REST(method=PUT, path=/notes/:key)
func (h *Handler) Save(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(err, errors.BadRequestf("the body is not a note"))
	}
	if err := h.svc.Save(c.Request().Context(), &entity.Note{Key: c.Param("key"), Text: req.Text}); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// Find answers the note of the key of the path.
// @REST(method=GET, path=/notes/:key)
func (h *Handler) Find(c echo.Context) error {
	note, err := h.svc.Find(c.Request().Context(), c.Param("key"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{Key: note.Key, Text: note.Text})
}
