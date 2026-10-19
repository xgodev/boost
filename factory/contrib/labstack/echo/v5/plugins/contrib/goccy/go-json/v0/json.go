package json

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	le "github.com/labstack/echo/v5"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new goccy/go-json plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewJson()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Json represents goccy/go-json plugin for echo server.
type Json struct {
	options *Options
}

// NewJsonWithOptions returns a new goccy/go-json plugin with options.
func NewJsonWithOptions(options *Options) *Json {
	return &Json{options: options}
}

// NewJsonWithConfigPath returns a new goccy/go-json plugin with options from config path.
func NewJsonWithConfigPath(path string) (*Json, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewJsonWithOptions(o), nil
}

// NewJson returns a new goccy/go-json plugin with default options.
func NewJson() (*Json, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewJsonWithOptions(o), nil
}

// Register registers this goccy/go-json plugin for echo server.
func (i *Json) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling goccy/go-json encoder in echo")

	server.Instance().JSONSerializer = JSONSerializer{}

	logger.Debug("goccy/go-json encoder successfully enabled in echo")

	return nil
}

// JSONSerializer implements echo's JSONSerializer with goccy/go-json.
type JSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d JSONSerializer) Serialize(c *le.Context, i any, indent string) error {
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d JSONSerializer) Deserialize(c *le.Context, i any) error {
	err := json.NewDecoder(c.Request().Body).Decode(i)

	var ute *json.UnmarshalTypeError
	var se *json.SyntaxError
	switch {
	case errors.As(err, &ute):
		return le.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset))
	case errors.As(err, &se):
		return le.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()))
	}
	return err
}
