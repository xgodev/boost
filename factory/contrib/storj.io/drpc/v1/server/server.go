package server

import (
	"context"
	"net"
	"strconv"
	"sync"

	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/wrapper/log"
	"storj.io/drpc"
	"storj.io/drpc/drpcmux"
	"storj.io/drpc/drpcserver"
)

// Plugin defines a drpc server Plugin function to execute.
type Plugin = factory.Plugin[*Server]

// Server represents a drpc server.
type Server struct {
	server  *drpcserver.Server
	mux     *drpcmux.Mux
	options *Options

	// mu guards stopped and cancel: Shutdown cancels the context Serve
	// runs with, and a Shutdown before Serve makes Serve return at once.
	mu      sync.Mutex
	stopped bool
	cancel  context.CancelFunc
	done    chan struct{}

	// listening is closed once Serve has bound its listener; addr is the
	// bound address, read only after listening is closed.
	listening  chan struct{}
	listenOnce sync.Once
	addr       net.Addr
}

// NewServer returns a new drpc server with default options.
func NewServer(ctx context.Context, plugins ...Plugin) (*Server, error) {
	opt, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewServerWithOptions(ctx, opt, plugins...)
}

// NewServerWithConfigPath returns a new drpc server with options from config path.
func NewServerWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*Server, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewServerWithOptions(ctx, options, plugins...)
}

// NewServerWithOptions returns a new drpc server with options. A plugin that
// fails stops the build: the error is returned and later plugins do not run.
func NewServerWithOptions(ctx context.Context, opt *Options, plugins ...Plugin) (*Server, error) {
	logger := log.FromContext(ctx)

	m := drpcmux.New()
	srv := &Server{
		mux: m,
		server: drpcserver.NewWithOptions(m, drpcserver.Options{
			CollectStats: opt.CollectStats,
			Log:          func(err error) { logger.Warn(err) },
		}),
		options:   opt,
		done:      make(chan struct{}),
		listening: make(chan struct{}),
	}

	if err := factory.ApplyPlugins(ctx, srv, plugins...); err != nil {
		return nil, err
	}

	return srv, nil
}

// Server returns the wrapped drpc server instance.
func (s *Server) Server() *drpcserver.Server {
	return s.server
}

// Mux returns the wrapped drpc mux instance.
func (s *Server) Mux() *drpcmux.Mux {
	return s.mux
}

// ServerOptions returns drpc server options.
func (s *Server) ServerOptions() *Options {
	return s.options
}

// Register dispatches the rpcs desc describes to impl.
func (s *Server) Register(impl interface{}, desc drpc.Description) error {
	return s.mux.Register(impl, desc)
}

// Serve starts the drpc server and blocks until Shutdown or ctx stops it.
// Errors are logged.
func (s *Server) Serve(ctx context.Context) {
	logger := log.FromContext(ctx)

	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.mu.Unlock()
	defer close(s.done)
	defer cancel()

	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", ":"+strconv.Itoa(s.options.Port))
	if err != nil {
		logger.Error(err)
		return
	}
	s.listenOnce.Do(func() {
		s.addr = ln.Addr()
		close(s.listening)
	})
	logger.Infof("drpc server started on %s", s.addr)

	if err := s.server.Serve(ctx, ln); err != nil {
		logger.Error(err)
	}
}

// Shutdown stops the drpc server, closing its listener and the connections
// it serves, and waits for Serve to return, bounded by ctx. Called before
// Serve, it makes a later Serve return at once.
func (s *Server) Shutdown(ctx context.Context) {
	s.mu.Lock()
	s.stopped = true
	cancel := s.cancel
	s.mu.Unlock()

	if cancel == nil {
		return
	}
	cancel()

	select {
	case <-s.done:
	case <-ctx.Done():
		log.FromContext(ctx).Warn(ctx.Err())
	}
}
