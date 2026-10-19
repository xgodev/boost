package fiber

import (
	"context"
	"net"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/wrapper/log"
)

// Plugin defines a fiber server Plugin function to execute.
type Plugin = factory.Plugin[*Server]

// Server represents a fiber v2 server.
type Server struct {
	app     *fiber.App
	options *Options
	config  fiber.Config

	// building is true while the plugins run: fiber reads its config only
	// in fiber.New, so the app is created after them. Their middlewares
	// and routes wait in middlewares and routes, and the middlewares are
	// added first, so they wrap every route.
	building    bool
	middlewares []func()
	routes      []func()

	// mu guards stopped and ln: Shutdown may run before, during or after
	// the listener is bound.
	mu      sync.Mutex
	stopped bool
	ln      net.Listener
	// listening is closed once Serve has bound its listener; addr is the
	// bound address, read only after listening is closed.
	listening  chan struct{}
	listenOnce sync.Once
	addr       net.Addr
}

// NewServer returns a new fiber server with default options.
func NewServer(ctx context.Context, plugins ...Plugin) (*Server, error) {
	opt, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewServerWithOptions(ctx, opt, plugins...)
}

// NewServerWithConfigPath returns a new fiber server with options from config path.
func NewServerWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*Server, error) {
	opt, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewServerWithOptions(ctx, opt, plugins...)
}

// NewServerWithOptions returns a new fiber server with options. The plugins
// may change the fiber config through Config; the app is created after all
// of them, with their middlewares before their routes. A plugin that fails
// stops the build: the error is returned and later plugins do not run.
func NewServerWithOptions(ctx context.Context, opt *Options, plugins ...Plugin) (*Server, error) {
	srv := &Server{
		options:   opt,
		config:    opt.Config.fiberConfig(),
		listening: make(chan struct{}),
		building:  true,
	}

	if err := factory.ApplyPlugins(ctx, srv, plugins...); err != nil {
		return nil, err
	}

	srv.app = fiber.New(srv.config)
	srv.building = false
	for _, add := range append(srv.middlewares, srv.routes...) {
		add()
	}
	srv.middlewares, srv.routes = nil, nil

	return srv, nil
}

// App returns the wrapped fiber app. It is nil while the plugins run.
func (s *Server) App() *fiber.App {
	return s.app
}

// ServerOptions returns fiber server options.
func (s *Server) ServerOptions() *Options {
	return s.options
}

// Config returns the fiber config the app is created with. Changes made
// after the plugins run have no effect.
func (s *Server) Config() *fiber.Config {
	return &s.config
}

// route adds a route now, or after the plugins while they run.
func (s *Server) route(add func()) {
	if s.building {
		s.routes = append(s.routes, add)
		return
	}
	add()
}

// Use adds a middleware, as fiber.App.Use. Outside the plugins, it only
// wraps the routes added after it, as in fiber.
func (s *Server) Use(args ...interface{}) {
	add := func() { s.app.Use(args...) }
	if s.building {
		s.middlewares = append(s.middlewares, add)
		return
	}
	add()
}

// Get adds the route path that matches a GET http method.
func (s *Server) Get(path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.Get(path, handlers...) })
}

// Head adds the route path that matches a HEAD http method.
func (s *Server) Head(path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.Head(path, handlers...) })
}

// Post adds the route path that matches a POST http method.
func (s *Server) Post(path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.Post(path, handlers...) })
}

// Put adds the route path that matches a PUT http method.
func (s *Server) Put(path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.Put(path, handlers...) })
}

// Delete adds the route path that matches a DELETE http method.
func (s *Server) Delete(path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.Delete(path, handlers...) })
}

// Connect adds the route path that matches a CONNECT http method.
func (s *Server) Connect(path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.Connect(path, handlers...) })
}

// Options adds the route path that matches an OPTIONS http method.
func (s *Server) Options(path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.Options(path, handlers...) })
}

// Trace adds the route path that matches a TRACE http method.
func (s *Server) Trace(path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.Trace(path, handlers...) })
}

// Patch adds the route path that matches a PATCH http method.
func (s *Server) Patch(path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.Patch(path, handlers...) })
}

// Add adds the route path that matches method.
func (s *Server) Add(method, path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.Add(method, path, handlers...) })
}

// All adds the route path that matches any http method.
func (s *Server) All(path string, handlers ...fiber.Handler) {
	s.route(func() { s.app.All(path, handlers...) })
}

// Static serves the files under root along prefix.
func (s *Server) Static(prefix, root string, config ...fiber.Static) {
	s.route(func() { s.app.Static(prefix, root, config...) })
}

// Route adds the routes fn builds along prefix.
func (s *Server) Route(prefix string, fn func(router fiber.Router), name ...string) {
	s.route(func() { s.app.Route(prefix, fn, name...) })
}

// Mount attaches the routes of app along prefix.
func (s *Server) Mount(prefix string, app *fiber.App) {
	s.route(func() { s.app.Mount(prefix, app) })
}

// Serve starts the fiber server and blocks until Shutdown stops it. Errors
// are logged.
func (s *Server) Serve(ctx context.Context) {
	logger := log.FromContext(ctx)

	if err := s.serve(ctx); err != nil {
		logger.Error(err)
	}
}

func (s *Server) serve(ctx context.Context) error {
	logger := log.FromContext(ctx)
	addr := ":" + strconv.Itoa(s.options.Port)
	// The app config has the defaults fiber.New filled in.
	cfg := s.app.Config()

	// Prefork binds in each child process, so fiber binds itself.
	if cfg.Prefork {
		logger.Infof("starting fiber server with prefork on %s", addr)
		return s.app.Listen(addr)
	}

	ln, err := (&net.ListenConfig{}).Listen(ctx, cfg.Network, addr)
	if err != nil {
		return err
	}

	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return ln.Close()
	}
	s.ln = ln
	s.mu.Unlock()

	s.listenOnce.Do(func() {
		s.addr = ln.Addr()
		close(s.listening)
	})
	logger.Infof("starting fiber server on %s", s.addr)

	return s.app.Listener(ln)
}

// Shutdown stops the fiber server gracefully, bounded by ctx. Called before
// Serve, it makes a later Serve return at once.
func (s *Server) Shutdown(ctx context.Context) {
	logger := log.FromContext(ctx)

	s.mu.Lock()
	s.stopped = true
	ln := s.ln
	s.mu.Unlock()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		logger.Warn(err)
	}
	// fiber closes the listener it serves; closing it again unblocks a
	// Serve that had bound it but not yet started.
	if ln != nil {
		_ = ln.Close() // intentional: fiber may have closed it already
	}
}
