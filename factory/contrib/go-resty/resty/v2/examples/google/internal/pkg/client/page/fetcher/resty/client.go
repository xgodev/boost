// Package resty fetches the page over HTTP, on the resty client of boost's
// factory (boost.factory.resty.clients.google).
package resty

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/examples/google/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/examples/google/internal/pkg/domain/service/page/fetcher"
	"github.com/xgodev/boost/model/errors"
)

type client struct {
	http *resty.Client
}

// New provides the fetcher on the client the generator builds with the
// health plugin, which registers a checker of the host with the keys under
// boost.factory.resty.clients.google.plugins.health.
// @Provide(index=0)
// @Inject(index=0, name=google, plugins=[health])
func New(c *resty.Client) fetcher.Service {
	return &client{http: c}
}

func (c *client) Fetch(ctx context.Context) (*entity.Page, error) {
	resp, err := c.http.R().SetContext(ctx).EnableTrace().Get("/")
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error fetching %s", c.http.BaseURL))
	}
	return &entity.Page{URL: c.http.BaseURL, Status: resp.StatusCode()}, nil
}
