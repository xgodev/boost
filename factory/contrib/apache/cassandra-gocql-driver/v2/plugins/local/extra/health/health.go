package health

import (
	"context"

	"github.com/apache/cassandra-gocql-driver/v2"
	"github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/wrapper/log"
)

// Register adds session to the health checks with the options in config; a
// config that does not load is returned as an error.
func Register(ctx context.Context, session *gocql.Session) error {
	h, err := NewHealth()
	if err != nil {
		return err
	}
	return h.Register(ctx, session)
}

// Health is the health plugin of the gocql v2 factory.
type Health struct {
	options *Options
}

// NewHealthWithOptions returns a health plugin with options.
func NewHealthWithOptions(options *Options) *Health {
	return &Health{options: options}
}

// NewHealthWithConfigPath returns a health plugin with the options of path.
func NewHealthWithConfigPath(path string) (*Health, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// NewHealth returns a health plugin with the default options.
func NewHealth() (*Health, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// Register adds session to the health checks.
func (i *Health) Register(ctx context.Context, session *gocql.Session) error {
	logger := log.FromContext(ctx).WithTypeOf(*i)

	logger.Trace("integrating gocql in health")

	checker := NewChecker(session)
	hc := health.NewHealthChecker(i.options.Name, i.options.Description, checker, i.options.Required, i.options.Enabled)
	health.Add(hc)

	logger.Debug("gocql successfully integrated in health")

	return nil
}
