package echo

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/factory/contrib/golang.org/x/net/v0/http2/server"
	"github.com/xgodev/boost/wrapper/log"
	"golang.org/x/net/http2/h2c"
)

// Plugin defines an echo server Plugin function to execute.
type Plugin = factory.Plugin[*Server]

// readTimeout bounds how long a request may take to be read, the same
// default echo v5 uses for the servers it starts.
const readTimeout = 30 * time.Second

// Server represents an echo v5 server.
type Server struct {
	instance *echo.Echo
	options  *Options
	h2cPath  string

	// once guards http, built on the first Serve or Shutdown.
	once sync.Once
	http *http.Server
	// listening is closed once Serve has bound its listener; addr is the
	// bound address, read only after listening is closed.
	listening  chan struct{}
	listenOnce sync.Once
	addr       net.Addr
}

// NewServer returns a new echo server with default options.
func NewServer(ctx context.Context, plugins ...Plugin) (*Server, error) {
	opt, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewServerWithOptions(ctx, opt, plugins...)
}

// NewServerWithConfigPath returns a new echo server with options from config path.
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

// NewServerWithOptions returns a new echo server with options. A plugin that
// fails stops the build: the error is returned and later plugins do not run.
func NewServerWithOptions(ctx context.Context, opt *Options, plugins ...Plugin) (*Server, error) {
	instance := echo.New()
	instance.Logger = WrapLogger(log.GetLogger())

	srv := &Server{
		instance:  instance,
		options:   opt,
		h2cPath:   root + hc2Root,
		listening: make(chan struct{}),
	}

	if err := factory.ApplyPlugins(ctx, srv, plugins...); err != nil {
		return nil, err
	}

	return srv, nil
}

// Instance returns the wrapped echo server instance.
func (s *Server) Instance() *echo.Echo {
	return s.instance
}

// Options returns echo server options.
func (s *Server) Options() *Options {
	return s.options
}

// GET registers the handler and middlewares for GET route at path.
func (s *Server) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.GET(path, h, m...)
}

// POST registers the handler and middlewares for POST route at path.
func (s *Server) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.POST(path, h, m...)
}

// PUT registers the handler and middlewares for PUT route at path.
func (s *Server) PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.PUT(path, h, m...)
}

// DELETE registers the handler and middlewares for DELETE route at path.
func (s *Server) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.DELETE(path, h, m...)
}

// OPTIONS registers the handler and middlewares for OPTIONS route at path.
func (s *Server) OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.OPTIONS(path, h, m...)
}

// PATCH registers the handler and middlewares for PATCH route at path.
func (s *Server) PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.PATCH(path, h, m...)
}

// HEAD registers the handler and middlewares for HEAD route at path.
func (s *Server) HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.HEAD(path, h, m...)
}

// CONNECT registers the handler and middlewares for CONNECT route at path.
func (s *Server) CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.CONNECT(path, h, m...)
}

// Group creates a router group with prefix and registers middlewares.
func (s *Server) Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group {
	return s.instance.Group(prefix, m...)
}

// Add registers the handler and middlewares for method route at path.
func (s *Server) Add(method string, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.Add(method, path, h, m...)
}

// Any registers the handler and middlewares for all method routes at path.
func (s *Server) Any(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.Any(path, h, m...)
}

// Use registers middlewares that will run after router.
func (s *Server) Use(middleware ...echo.MiddlewareFunc) {
	s.instance.Use(middleware...)
}

// Static registers a static path to server static files at root dir.
func (s *Server) Static(prefix, root string, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.Static(prefix, root, m...)
}

// Match registers the handler and middlewares for method routes at path.
func (s *Server) Match(methods []string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) echo.Routes {
	return s.instance.Match(methods, path, handler, middleware...)
}

// Pre registers middlewares that will run before router.
func (s *Server) Pre(middleware ...echo.MiddlewareFunc) {
	s.instance.Pre(middleware...)
}

// File registers a path to a server file.
func (s *Server) File(path, file string, m ...echo.MiddlewareFunc) echo.RouteInfo {
	return s.instance.File(path, file, m...)
}

// Serve starts the echo server and blocks until Shutdown stops it. Errors are
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
		srv.Handler = h2c.NewHandler(s.instance, h2s)
		logger.Infof("starting Echo H2C Server. https://echo.labstack.com/")
	case s.options.TLS.Enabled:
		cfg, err := s.tlsConfig()
		if err != nil {
			return err
		}
		srv.TLSConfig = cfg
		logger.Infof("starting Echo TLS Server. https://echo.labstack.com/")
	default:
		logger.Infof("starting Echo Server. https://echo.labstack.com/")
	}

	if !s.options.HideBanner {
		logger.Infof("Echo v%s", echo.Version)
	}

	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", ":"+strconv.Itoa(s.options.Port))
	if err != nil {
		return err
	}
	s.listenOnce.Do(func() {
		s.addr = ln.Addr()
		close(s.listening)
	})
	logger.Debugf("echo server listening on %s", s.addr)

	if srv.TLSConfig != nil {
		return srv.ServeTLS(ln, "", "")
	}
	return srv.Serve(ln)
}

// httpServer builds, once, the http.Server Serve runs and Shutdown stops.
func (s *Server) httpServer() *http.Server {
	s.once.Do(func() {
		s.http = &http.Server{
			Handler:     s.instance,
			ErrorLog:    slog.NewLogLogger(s.instance.Logger.Handler(), slog.LevelError),
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

// Shutdown stops the echo server gracefully, bounded by ctx. Called before
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
			Organization: []string{"echo http server"},
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
