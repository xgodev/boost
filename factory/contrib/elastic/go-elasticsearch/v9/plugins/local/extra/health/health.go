package health

import (
	"context"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/wrapper/log"
)

// Health represents elasticsearch health.
type Health struct {
	options *Options
}

// NewHealthWithOptions returns a health with the options provided.
func NewHealthWithOptions(options *Options) *Health {
	return &Health{options: options}
}

// NewHealthWithConfigPath returns a health with options from config path.
func NewHealthWithConfigPath(path string) (*Health, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// NewHealth returns a health with default options.
func NewHealth() (*Health, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewHealthWithOptions(o), nil
}

// Register registers a new checker in the health package.
func (i *Health) Register(ctx context.Context, client *elasticsearch.Client) error {

	logger := log.FromContext(ctx).WithTypeOf(*i)

	logger.Trace("integrating elasticsearch in health")

	checker := NewChecker(client)
	hc := health.NewHealthChecker(i.options.Name, i.options.Description, checker, i.options.Required, i.options.Enabled)
	health.Add(hc)

	logger.Debug("elasticsearch successfully integrated in health")

	return nil
}

// Register is a factory plugin that registers client in the health package
// with the default options.
func Register(ctx context.Context, client *elasticsearch.Client) error {
	h, err := NewHealth()
	if err != nil {
		return err
	}
	return h.Register(ctx, client)
}
