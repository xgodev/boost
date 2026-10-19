package multiserver

import (
	"context"
	"net/http"

	"github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new multi server check route plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	p, err := NewMultiServer()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// MultiServer represents multi server check route plugin for chi server.
type MultiServer struct {
	options *Options
}

// NewMultiServerWithOptions returns a new multi server check route plugin with options.
func NewMultiServerWithOptions(options *Options) *MultiServer {
	return &MultiServer{options: options}
}

// NewMultiServerWithConfigPath returns a new multi server check route plugin with options from config path.
func NewMultiServerWithConfigPath(path string) (*MultiServer, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewMultiServerWithOptions(o), nil
}

// NewMultiServer returns a new multi server check route plugin with default options.
func NewMultiServer() (*MultiServer, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewMultiServerWithOptions(o), nil
}

// Register registers this multi server check route plugin for chi server.
func (i *MultiServer) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling multi server check route in chi")

	route := i.options.Route
	server.Get(route, handler)
	logger.Debugf("multi server check route configured on %s in chi", route)

	logger.Debug("multi server check route successfully enabled in chi")

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := multiserver.Check(r.Context()); err != nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}
	_, _ = w.Write([]byte("OK")) // intentional: the caller is gone when the write fails
}
