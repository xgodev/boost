package log

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new request log plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	h, err := NewLog()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Log represents request log plugin for chi server.
type Log struct {
	options *Options
}

// NewLogWithOptions returns a new request log plugin with options.
func NewLogWithOptions(options *Options) *Log {
	return &Log{options: options}
}

// NewLogWithConfigPath returns a new request log plugin with options from config path.
func NewLogWithConfigPath(path string) (*Log, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewLogWithOptions(o), nil
}

// NewLog returns a new request log plugin with default options.
func NewLog() (*Log, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewLogWithOptions(o), nil
}

// Register registers this request log plugin for chi server.
func (i *Log) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling request log middleware in chi")

	server.Use(loggerMiddleware(i.options.Level))

	logger.Debug("request log middleware successfully enabled in chi")

	return nil
}

// loggerMiddleware logs every request through the boost logger, at level, and
// puts that logger in the request context for the handlers.
func loggerMiddleware(level string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logger := log.FromContext(r.Context()).WithField("requestUri", r.RequestURI)
			r = r.WithContext(logger.ToContext(r.Context()))

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			// A handler that writes nothing answers 200.
			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}

			reqSize := r.Header.Get("Content-Length")
			if reqSize == "" {
				reqSize = "0"
			}

			method := logger.Debugf
			switch level {
			case "TRACE":
				method = logger.Tracef
			case "INFO":
				method = logger.Infof
			}

			method("%s %s %-7s %s %3d %s %s %13v %s %s",
				r.RemoteAddr,
				r.Host,
				r.Method,
				r.RequestURI,
				status,
				reqSize,
				strconv.Itoa(ww.BytesWritten()),
				time.Since(start).String(),
				r.Referer(),
				r.UserAgent(),
			)
		})
	}
}
