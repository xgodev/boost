package log

import (
	"context"
	"time"

	f "github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new request log middleware plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewLog()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Log represents request log middleware plugin for fiber server.
type Log struct {
	options *Options
}

// NewLogWithOptions returns a new request log middleware plugin with options.
func NewLogWithOptions(options *Options) *Log {
	return &Log{options: options}
}

// NewLogWithConfigPath returns a new request log middleware plugin with options from config path.
func NewLogWithConfigPath(path string) (*Log, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewLogWithOptions(o), nil
}

// NewLog returns a new request log middleware plugin with default options.
func NewLog() (*Log, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewLogWithOptions(o), nil
}

// Register registers this request log middleware plugin for fiber server.
func (i *Log) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling request log middleware in fiber")

	server.Use(loggerMiddleware(i.options.Level))

	logger.Debug("request log middleware successfully enabled in fiber")

	return nil
}

// loggerMiddleware logs every request through the boost logger, at level,
// and puts that logger in the user context for the handlers. An error from
// the handlers goes through the app error handler first, so the logged
// status is the one the client gets.
func loggerMiddleware(level string) f.Handler {
	return func(c *f.Ctx) error {
		start := time.Now()
		logger := log.FromContext(c.UserContext()).WithField("requestUri", c.OriginalURL())
		c.SetUserContext(logger.ToContext(c.UserContext()))

		if err := c.Next(); err != nil {
			if herr := c.App().ErrorHandler(c, err); herr != nil {
				_ = c.SendStatus(f.StatusInternalServerError) // intentional: the caller is gone when the write fails
			}
		}

		method := logger.Debugf
		switch level {
		case "TRACE":
			method = logger.Tracef
		case "INFO":
			method = logger.Infof
		}

		method("%s %s %-7s %s %3d %d %d %13v %s %s",
			c.IP(),
			c.Hostname(),
			c.Method(),
			c.OriginalURL(),
			c.Response().StatusCode(),
			len(c.Request().Body()),
			len(c.Response().Body()),
			time.Since(start).String(),
			c.Get(f.HeaderReferer),
			c.Get(f.HeaderUserAgent),
		)
		return nil
	}
}
