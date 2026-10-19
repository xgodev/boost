package json

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bytedance/sonic/decoder"
	"github.com/bytedance/sonic/encoder"
	le "github.com/labstack/echo/v5"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new bytedance/sonic plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewSonic()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Sonic represents bytedance/sonic plugin for echo server.
type Sonic struct {
	options *Options
}

// NewSonicWithOptions returns a new bytedance/sonic plugin with options.
func NewSonicWithOptions(options *Options) *Sonic {
	return &Sonic{options: options}
}

// NewSonicWithConfigPath returns a new bytedance/sonic plugin with options from config path.
func NewSonicWithConfigPath(path string) (*Sonic, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewSonicWithOptions(o), nil
}

// NewSonic returns a new bytedance/sonic plugin with default options.
func NewSonic() (*Sonic, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewSonicWithOptions(o), nil
}

// Register registers this bytedance/sonic plugin for echo server.
func (i *Sonic) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling bytedance/sonic encoder in echo")

	server.Instance().JSONSerializer = JSONSerializer{}

	logger.Debug("bytedance/sonic encoder successfully enabled in echo")

	return nil
}

// JSONSerializer implements echo's JSONSerializer with bytedance/sonic.
type JSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d JSONSerializer) Serialize(c *le.Context, i any, indent string) error {
	enc := encoder.NewStreamEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d JSONSerializer) Deserialize(c *le.Context, i any) error {
	err := decoder.NewStreamDecoder(c.Request().Body).Decode(i)

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
