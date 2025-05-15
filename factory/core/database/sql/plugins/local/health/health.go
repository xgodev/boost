package health

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/wrapper/log"
)

// HealthPlugin implements the factory’s Plugin hooks for health checks.
type HealthPlugin struct {
	options *Options
}

// NewHealthPluginWithOptions constructs the plugin from explicit Options.
func NewHealthPluginWithOptions(opts *Options) *HealthPlugin {
	return &HealthPlugin{options: opts}
}

// NewHealthPluginWithConfigPath loads Options from a config file.
func NewHealthPluginWithConfigPath(path string) (*HealthPlugin, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewHealthPluginWithOptions(opts), nil
}

// NewHealthPlugin constructs the plugin with default Options.
func NewHealthPlugin() *HealthPlugin {
	opts, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return NewHealthPluginWithOptions(opts)
}

// WrapConnector is a no-op: health checks don’t need to wrap the driver.
func (h *HealthPlugin) WrapConnector(ctx context.Context, connector driver.Connector) (driver.Connector, error) {
	return connector, nil
}

// InitDB is called immediately after sql.OpenDB: it registers the health check.
func (h *HealthPlugin) InitDB(ctx context.Context, db *sql.DB) error {
	logger := log.FromContext(ctx).WithTypeOf(*h)
	logger.Trace("integrating sql in health")

	checker := NewChecker(db)
	hc := health.NewHealthChecker(
		h.options.Name,
		h.options.Description,
		checker,
		h.options.Required,
		h.options.Enabled,
	)
	health.Add(hc)

	logger.Debug("sql successfully integrated in health")
	return nil
}
