package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Plugin customizes the server right after it is built, before Serve.
type Plugin func(context.Context, *Server) error

// Server wraps an *mcp.Server and serves it over the configured transport.
type Server struct {
	instance *mcp.Server
	options  *Options

	// transport is what Serve runs STDIO on; tests swap in an in-memory pair.
	transport mcp.Transport

	mu     sync.Mutex
	http   *http.Server
	cancel context.CancelFunc
}

// NewServer builds a server from the default config path.
func NewServer(ctx context.Context, plugins ...Plugin) (*Server, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewServerWithOptions(ctx, o, plugins...)
}

// NewServerWithConfigPath builds a server from the given config path.
func NewServerWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*Server, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewServerWithOptions(ctx, o, plugins...)
}

// NewServerWithOptions builds a server from o and runs the plugins on it.
func NewServerWithOptions(ctx context.Context, o *Options, plugins ...Plugin) (*Server, error) {
	o.Transport = strings.ToUpper(o.Transport)
	if o.Transport != TransportHTTP && o.Transport != TransportSTDIO {
		return nil, errors.NotValidf("mcp server transport %q: want %s or %s", o.Transport, TransportHTTP, TransportSTDIO)
	}

	srv := &Server{
		instance: mcp.NewServer(
			&mcp.Implementation{Name: o.Name, Version: o.Version},
			&mcp.ServerOptions{Instructions: o.Instructions, KeepAlive: o.KeepAlive},
		),
		options:   o,
		transport: &mcp.StdioTransport{},
	}

	for _, plugin := range plugins {
		if err := plugin(ctx, srv); err != nil {
			return nil, err
		}
	}

	return srv, nil
}

// Instance returns the wrapped server; register tools, prompts and resources on it.
func (s *Server) Instance() *mcp.Server {
	return s.instance
}

// Options returns the options the server was built with.
func (s *Server) Options() *Options {
	return s.options
}

// Handler returns the streamable HTTP handler, to mount the server on an
// existing HTTP server (echo, net/http) instead of its own port.
func (s *Server) Handler() http.Handler {
	return mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return s.instance },
		&mcp.StreamableHTTPOptions{
			Stateless:      s.options.Stateless,
			JSONResponse:   s.options.JSONResponse,
			SessionTimeout: s.options.SessionTimeout,
		},
	)
}

// Serve blocks serving the configured transport until Shutdown or a serve error.
func (s *Server) Serve(ctx context.Context) {
	if s.options.Transport == TransportSTDIO {
		s.serveSTDIO(ctx)
		return
	}
	s.serveHTTP(ctx)
}

func (s *Server) serveHTTP(ctx context.Context) {
	logger := log.FromContext(ctx)

	mux := http.NewServeMux()
	mux.Handle(s.options.Path, s.Handler())
	hs := &http.Server{Handler: mux}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.options.Port))
	if err != nil {
		logger.Errorf("mcp server: listen on port %d: %v", s.options.Port, err)
		return
	}

	s.mu.Lock()
	s.http = hs
	s.mu.Unlock()

	logger.Infof("mcp server listening on %s%s", l.Addr(), s.options.Path)
	if err := hs.Serve(l); err != nil && err != http.ErrServerClosed {
		logger.Errorf("mcp server: %v", err)
	}
}

func (s *Server) serveSTDIO(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.mu.Lock()
	s.cancel = cancel
	s.mu.Unlock()

	if err := s.instance.Run(ctx, s.transport); err != nil && ctx.Err() == nil {
		log.FromContext(ctx).Errorf("mcp server: %v", err)
	}
}

// Shutdown stops a running Serve; it is a no-op when nothing is serving.
func (s *Server) Shutdown(ctx context.Context) {
	s.mu.Lock()
	hs, cancel := s.http, s.cancel
	s.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if hs != nil {
		if err := hs.Shutdown(ctx); err != nil {
			log.FromContext(ctx).Errorf("mcp server: shutdown: %v", err)
		}
	}
}
