// Package resty fetches the pages over HTTP, on resty clients of boost's
// factory: one per site, each with its own keys under
// boost.factory.resty.clients.
package resty

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/examples/cobra/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/examples/cobra/internal/pkg/domain/service/page/fetcher"
	"github.com/xgodev/boost/model/errors"
)

type client struct {
	http *resty.Client
}

// NewGoogle provides the google fetcher on the google client, which the
// generator builds without plugins.
// @Provide(name=google, index=0)
// @Inject(index=0, name=google)
func NewGoogle(c *resty.Client) fetcher.Service {
	return &client{http: c}
}

// NewBing provides the bing fetcher on the bing client, which the generator
// builds with the log plugin.
// @Provide(name=bing, index=0)
// @Inject(index=0, name=bing, plugins=[log])
func NewBing(c *resty.Client) fetcher.Service {
	return &client{http: c}
}

func (c *client) Fetch(ctx context.Context) (*entity.Page, error) {
	resp, err := c.http.R().SetContext(ctx).Get("/")
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error fetching %s", c.http.BaseURL))
	}
	return &entity.Page{URL: c.http.BaseURL, Status: resp.StatusCode()}, nil
}
