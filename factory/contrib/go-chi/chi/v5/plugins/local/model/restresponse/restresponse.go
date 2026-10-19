package restresponse

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	response "github.com/xgodev/boost/model/restresponse"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new resource status route plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	p, err := NewStatus()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Status represents resource status route plugin for chi server.
type Status struct {
	options *Options
}

// NewStatusWithOptions returns a new resource status route plugin with options.
func NewStatusWithOptions(options *Options) *Status {
	return &Status{options: options}
}

// NewStatusWithConfigPath returns a new resource status route plugin with options from config path.
func NewStatusWithConfigPath(path string) (*Status, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewStatusWithOptions(o), nil
}

// NewStatus returns a new resource status route plugin with default options.
func NewStatus() (*Status, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewStatusWithOptions(o), nil
}

// Register registers this resource status route plugin for chi server.
func (i *Status) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling resource status route in chi")

	route := i.options.Route
	server.Get(route, handler)
	logger.Debugf("resource status route configured on %s in chi", route)

	logger.Debug("resource status route successfully enabled in chi")

	return nil
}

func handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response.NewResourceStatus()) // intentional: the caller is gone when the write fails
}
