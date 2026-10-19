// Package keeper serves the readings: PUT and GET /readings/:sensor.
package keeper

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache-binary/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/wrapper/cache/examples/freecache-binary/internal/pkg/domain/service/reading/keeper"
)

// Body is the body of PUT and GET /readings/:sensor.
type Body struct {
	Location    uint16 `json:"location"`
	Timestamp   uint32 `json:"timestamp"`
	Temperature int16  `json:"temperature"`
}

// Handler keeps readings.
type Handler struct {
	svc keeper.Service
}

// New provides the handler to the generated rest door.
// @Provide(index=0)
func New(svc keeper.Service) *Handler {
	return &Handler{svc: svc}
}

func sensor(c echo.Context) (uint32, error) {
	s, err := strconv.ParseUint(c.Param("sensor"), 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, errors.BadRequestf("sensor %q is not a number", c.Param("sensor")))
	}
	return uint32(s), nil
}

// Save keeps the reading of the body as the last one of the sensor of the
// path.
// @REST(method=PUT, path=/readings/:sensor)
func (h *Handler) Save(c echo.Context) error {
	s, err := sensor(c)
	if err != nil {
		return err
	}
	var b Body
	if err := c.Bind(&b); err != nil {
		return errors.Wrap(err, errors.BadRequestf("the body is not a reading"))
	}
	r := &entity.Reading{Sensor: s, Location: b.Location, Timestamp: b.Timestamp, Temperature: b.Temperature}
	if err := h.svc.Save(c.Request().Context(), r); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// Find answers the last reading of the sensor of the path.
// @REST(method=GET, path=/readings/:sensor)
func (h *Handler) Find(c echo.Context) error {
	s, err := sensor(c)
	if err != nil {
		return err
	}
	r, err := h.svc.Find(c.Request().Context(), s)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Body{Location: r.Location, Timestamp: r.Timestamp, Temperature: r.Temperature})
}
