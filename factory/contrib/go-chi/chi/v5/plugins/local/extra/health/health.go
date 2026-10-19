package health

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	response "github.com/xgodev/boost/model/restresponse"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new health route plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	p, err := NewHealth()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Health represents health route plugin for chi server.
type Health struct {
	options *Options
}

// NewHealthWithOptions returns a new health route plugin with options.
func NewHealthWithOptions(options *Options) *Health {
	return &Health{options: options}
}

// NewHealthWithConfigPath returns a new health route plugin with options from config path.
func NewHealthWithConfigPath(path string) (*Health, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// NewHealth returns a new health route plugin with default options.
func NewHealth() (*Health, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// Register registers this health route plugin for chi server.
func (i *Health) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling health route in chi")

	route := i.options.Route
	server.Get(route, handler)
	logger.Debugf("health route configured on %s in chi", route)

	logger.Debug("health route successfully enabled in chi")

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp, httpCode := response.NewHealth(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_ = json.NewEncoder(w).Encode(resp) // intentional: the caller is gone when the write fails
}
