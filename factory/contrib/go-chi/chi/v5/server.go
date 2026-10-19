package chi

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	c "github.com/go-chi/chi/v5"
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/factory/contrib/golang.org/x/net/v0/http2/server"
	"github.com/xgodev/boost/wrapper/log"
	"golang.org/x/net/http2/h2c"
)

// Plugin defines a chi server Plugin function to execute.
type Plugin = factory.Plugin[*Server]

// readTimeout bounds how long a request may take to be read.
const readTimeout = 30 * time.Second

// Server represents a chi v5 server.
type Server struct {
	mux     *c.Mux
	options *Options
	h2cPath string

	// building is true while the plugins run: their routes wait in routes
	// and are added after all of them, as chi refuses a middleware after a
	// route.
	building bool
	routes   []func()

	// once guards http, built on the first Serve or Shutdown.
	once sync.Once
	http *http.Server
	// listening is closed once Serve has bound its listener; addr is the
	// bound address, read only after listening is closed.
	listening  chan struct{}
	listenOnce sync.Once
	addr       net.Addr
}

// NewServer returns a new chi server with default options.
func NewServer(ctx context.Context, plugins ...Plugin) (*Server, error) {
	opt, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewServerWithOptions(ctx, opt, plugins...)
}

// NewServerWithConfigPath returns a new chi server with options from config path.
func NewServerWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*Server, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	srv, err := NewServerWithOptions(ctx, options, plugins...)
	if err != nil {
		return nil, err
	}
	srv.h2cPath = path + hc2Root
	return srv, nil
}

// NewServerWithOptions returns a new chi server with options. The routes the
// plugins add are added after all of them, so every plugin middleware wraps
// them. A plugin that fails stops the build: the error is returned and later
// plugins do not run.
func NewServerWithOptions(ctx context.Context, opt *Options, plugins ...Plugin) (*Server, error) {
	srv := &Server{
		mux:       c.NewRouter(),
		options:   opt,
		h2cPath:   root + hc2Root,
		listening: make(chan struct{}),
		building:  true,
	}

	if err := factory.ApplyPlugins(ctx, srv, plugins...); err != nil {
		return nil, err
	}

	srv.building = false
	for _, route := range srv.routes {
		route()
	}
	srv.routes = nil

	return srv, nil
}

// Mux returns the wrapped chi router.
func (s *Server) Mux() *c.Mux {
	return s.mux
}

// ServerOptions returns chi server options.
func (s *Server) ServerOptions() *Options {
	return s.options
}

// route adds a route now, or after the plugins while they run.
func (s *Server) route(add func()) {
	if s.building {
		s.routes = append(s.routes, add)
		return
	}
	add()
}

// Use appends middlewares to the router stack. Outside the plugins, they
// must come before the first route, as in chi.
func (s *Server) Use(middlewares ...func(http.Handler) http.Handler) {
	s.mux.Use(middlewares...)
}

// With returns an inline router that runs the middlewares for the routes
// added to it. It fixes the middleware stack, so plugins must not call it.
func (s *Server) With(middlewares ...func(http.Handler) http.Handler) c.Router {
	return s.mux.With(middlewares...)
}

// Get adds the route pattern that matches a GET http method.
func (s *Server) Get(pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.Get(pattern, handlerFn) })
}

// Post adds the route pattern that matches a POST http method.
func (s *Server) Post(pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.Post(pattern, handlerFn) })
}

// Put adds the route pattern that matches a PUT http method.
func (s *Server) Put(pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.Put(pattern, handlerFn) })
}

// Patch adds the route pattern that matches a PATCH http method.
func (s *Server) Patch(pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.Patch(pattern, handlerFn) })
}

// Delete adds the route pattern that matches a DELETE http method.
func (s *Server) Delete(pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.Delete(pattern, handlerFn) })
}

// Head adds the route pattern that matches a HEAD http method.
func (s *Server) Head(pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.Head(pattern, handlerFn) })
}

// Options adds the route pattern that matches an OPTIONS http method.
func (s *Server) Options(pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.Options(pattern, handlerFn) })
}

// Connect adds the route pattern that matches a CONNECT http method.
func (s *Server) Connect(pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.Connect(pattern, handlerFn) })
}

// Trace adds the route pattern that matches a TRACE http method.
func (s *Server) Trace(pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.Trace(pattern, handlerFn) })
}

// Method adds the route pattern that matches method.
func (s *Server) Method(method, pattern string, handler http.Handler) {
	s.route(func() { s.mux.Method(method, pattern, handler) })
}

// MethodFunc adds the route pattern that matches method.
func (s *Server) MethodFunc(method, pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.MethodFunc(method, pattern, handlerFn) })
}

// Handle adds the route pattern that matches any http method.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.route(func() { s.mux.Handle(pattern, handler) })
}

// HandleFunc adds the route pattern that matches any http method.
func (s *Server) HandleFunc(pattern string, handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.HandleFunc(pattern, handlerFn) })
}

// Route mounts a sub-router, built by fn, along pattern.
func (s *Server) Route(pattern string, fn func(r c.Router)) {
	s.route(func() { s.mux.Route(pattern, fn) })
}

// Group adds the routes fn builds on an inline router with its own
// middleware stack.
func (s *Server) Group(fn func(r c.Router)) {
	s.route(func() { s.mux.Group(fn) })
}

// Mount attaches handler along pattern.
func (s *Server) Mount(pattern string, handler http.Handler) {
	s.route(func() { s.mux.Mount(pattern, handler) })
}

// NotFound sets the handler for routes that do not match.
func (s *Server) NotFound(handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.NotFound(handlerFn) })
}

// MethodNotAllowed sets the handler for a route without the request method.
func (s *Server) MethodNotAllowed(handlerFn http.HandlerFunc) {
	s.route(func() { s.mux.MethodNotAllowed(handlerFn) })
}

// Serve starts the chi server and blocks until Shutdown stops it. Errors are
// logged; the http.ErrServerClosed that Shutdown causes is not an error.
func (s *Server) Serve(ctx context.Context) {
	logger := log.FromContext(ctx)

	if err := s.serve(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err)
	}
}

func (s *Server) serve(ctx context.Context) error {
	logger := log.FromContext(ctx)

	srv := s.httpServer()

	switch {
	case s.options.Protocol == "H2C":
		h2s, err := server.NewServerWithPath(s.h2cPath)
		if err != nil {
			return err
		}
		srv.Handler = h2c.NewHandler(s.mux, h2s)
		logger.Infof("starting chi H2C server")
	case s.options.TLS.Enabled:
		cfg, err := s.tlsConfig()
		if err != nil {
			return err
		}
		srv.TLSConfig = cfg
		logger.Infof("starting chi TLS server")
	default:
		logger.Infof("starting chi server")
	}

	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", ":"+strconv.Itoa(s.options.Port))
	if err != nil {
		return err
	}
	s.listenOnce.Do(func() {
		s.addr = ln.Addr()
		close(s.listening)
	})
	logger.Debugf("chi server listening on %s", s.addr)

	if srv.TLSConfig != nil {
		return srv.ServeTLS(ln, "", "")
	}
	return srv.Serve(ln)
}

// httpServer builds, once, the http.Server Serve runs and Shutdown stops.
func (s *Server) httpServer() *http.Server {
	s.once.Do(func() {
		s.http = &http.Server{
			Handler:     s.mux,
			ReadTimeout: readTimeout,
		}
	})
	return s.http
}

func (s *Server) tlsConfig() (*tls.Config, error) {
	var (
		cert, key []byte
		err       error
	)

	if s.options.TLS.Type == "AUTO" {
		cert, key, err = generateCertificate(s.options.TLS.Auto.Host)
	} else {
		cert, err = os.ReadFile(s.options.TLS.File.Cert)
		if err == nil {
			key, err = os.ReadFile(s.options.TLS.File.Key)
		}
	}
	if err != nil {
		return nil, err
	}

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{pair}, MinVersion: tls.VersionTLS12}
	if s.options.DisableHTTP2 {
		cfg.NextProtos = []string{"http/1.1"}
	}
	return cfg, nil
}

// Shutdown stops the chi server gracefully, bounded by ctx. Called before
// Serve, it makes a later Serve return at once.
func (s *Server) Shutdown(ctx context.Context) {
	logger := log.FromContext(ctx)
	if err := s.httpServer().Shutdown(ctx); err != nil {
		logger.Warn(err)
	}
}

// generateCertificate generates a test certificate and private key based on the given host.
func generateCertificate(host string) ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"chi http server"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		DNSNames:              []string{host},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	crt := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})

	return crt, key, nil
}
