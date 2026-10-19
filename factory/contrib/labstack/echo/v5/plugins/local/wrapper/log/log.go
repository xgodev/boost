package log

import (
	"context"
	"strconv"
	"time"

	e "github.com/labstack/echo/v5"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new request log plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewLog()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Log represents request log plugin for echo server.
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

// Register registers this request log plugin for echo server.
func (i *Log) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling request log middleware in echo")

	server.Use(loggerMiddleware(i.options.Level))

	logger.Debug("request log middleware successfully enabled in echo")

	return nil
}

// loggerMiddleware logs every request through the boost logger, at level, and
// puts that logger in the request context for the handlers.
func loggerMiddleware(level string) e.MiddlewareFunc {
	return func(next e.HandlerFunc) e.HandlerFunc {
		return func(c *e.Context) error {
			req := c.Request()
			start := time.Now()

			logger := log.FromContext(req.Context()).WithField("requestUri", req.RequestURI)
			req = req.WithContext(logger.ToContext(req.Context()))
			c.SetRequest(req)

			// The error is handled here, as echo v4's c.Error did, so the
			// line logs the status the client gets.
			err := next(c)
			if err != nil {
				c.Echo().HTTPErrorHandler(c, err)
			}

			res, status := e.ResolveResponseStatus(c.Response(), err)
			var size int64
			if res != nil {
				size = res.Size
			}

			reqSize := req.Header.Get(e.HeaderContentLength)
			if reqSize == "" {
				reqSize = "0"
			}

			method := logger.Debugf
			if pol, ok := errors.IgnoreOf(err); !ok || pol&errors.IgnoreSilenceLog == 0 {
				switch level {
				case "TRACE":
					method = logger.Tracef
				case "INFO":
					method = logger.Infof
				}
			}

			method("%s %s %-7s %s %3d %s %s %13v %s %s",
				c.RealIP(),
				req.Host,
				req.Method,
				req.RequestURI,
				status,
				reqSize,
				strconv.FormatInt(size, 10),
				time.Since(start).String(),
				req.Referer(),
				req.UserAgent(),
			)

			return nil
		}
	}
}
