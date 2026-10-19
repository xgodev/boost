// Package resty fetches the page over HTTP, on the resty client of boost's
// factory (boost.factory.resty.clients.google).
package resty

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/googlecall/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/googlecall/internal/pkg/domain/service/page/fetcher"
	"github.com/xgodev/boost/model/errors"
)

type client struct {
	http *resty.Client
	url  string
}

// New provides the fetcher on the client the generator builds, with the log
// and requestid plugins.
// @Provide(index=0)
// @Inject(index=0, name=google, plugins=[log, requestid])
func New(c *resty.Client, o *Options) fetcher.Service {
	return &client{http: c, url: o.URL}
}

func (c *client) Fetch(ctx context.Context) (*entity.Page, error) {
	resp, err := c.http.R().SetContext(ctx).EnableTrace().Get(c.url)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error fetching %s", c.url))
	}
	return &entity.Page{URL: c.url, Status: resp.StatusCode()}, nil
}
