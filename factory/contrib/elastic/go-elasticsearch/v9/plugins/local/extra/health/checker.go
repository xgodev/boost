package health

import (
	"context"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xgodev/boost/model/errors"
)

// Checker represents elasticsearch health check.
type Checker struct {
	client *elasticsearch.Client
}

// Check pings the cluster; a ping answered with an error status is a
// ServiceUnavailable.
func (c *Checker) Check(ctx context.Context) error {
	res, err := c.client.Ping(c.client.Ping.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return errors.ServiceUnavailablef("elasticsearch answered the ping with %s", res.Status())
	}
	return nil
}

// NewChecker creates the elasticsearch health check.
func NewChecker(client *elasticsearch.Client) *Checker {
	return &Checker{client: client}
}
