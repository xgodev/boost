package requestid

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new request id middleware plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	p, err := NewRequestID()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// RequestID represents request id middleware plugin for chi server.
type RequestID struct {
	options *Options
}

// NewRequestIDWithOptions returns a new request id middleware plugin with options.
func NewRequestIDWithOptions(options *Options) *RequestID {
	return &RequestID{options: options}
}

// NewRequestIDWithConfigPath returns a new request id middleware plugin with options from config path.
func NewRequestIDWithConfigPath(path string) (*RequestID, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewRequestIDWithOptions(o), nil
}

// NewRequestID returns a new request id middleware plugin with default options.
func NewRequestID() (*RequestID, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewRequestIDWithOptions(o), nil
}

// Register registers this request id middleware plugin for chi server.
func (i *RequestID) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling request id middleware in chi")

	server.Use(middleware.RequestID, responseHeader)

	logger.Debug("request id middleware successfully enabled in chi")

	return nil
}

// responseHeader answers the request id, the caller's X-Request-Id or the one
// middleware.RequestID made, in the X-Request-Id header.
func responseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := middleware.GetReqID(r.Context()); id != "" {
			w.Header().Set(middleware.RequestIDHeader, id)
		}
		next.ServeHTTP(w, r)
	})
}
