// Package resty reserves stock by asking app api over HTTP
// (POST /stocks/:sku/reserve), chosen by
// INTERNAL_PKG_DOMAIN_SERVICE_STOCK_RESERVER_IMPLEMENTATION=resty. Its client
// is the resty client builder generates for the name stock, with its own
// config (boost.factory.resty.clients.stock.*).
package resty

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/reserver"
	"github.com/xgodev/boost/model/errors"
)

type client struct {
	http *resty.Client
}

// New provides the remote stock reserver; builder builds its client with
// the log, requestid and prometheus plugins.
// @Provide(index=0)
// @Inject(index=0, name=stock, plugins=[log, requestid, prometheus])
func New(c *resty.Client) reserver.Service {
	return &client{http: c}
}

// stock is the body app api answers with.
type stock struct {
	SKU   string `json:"sku"`
	Units int64  `json:"units"`
}

func (c *client) Reserve(ctx context.Context, sku string) (*entity.Stock, error) {
	var left stock
	resp, err := c.http.R().
		SetContext(ctx).
		SetPathParam("sku", sku).
		SetResult(&left).
		Post("/stocks/{sku}/reserve")
	if err != nil {
		return nil, errors.NewInternal(err, "asking app api to reserve "+sku)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		return &entity.Stock{SKU: left.SKU, Units: left.Units}, nil
	case http.StatusNotFound:
		return nil, errors.NotFoundf("%s is not in stock", sku)
	case http.StatusConflict:
		return nil, errors.Conflictf("%s is sold out", sku)
	default:
		return nil, errors.Internalf("app api answered %s reserving %s", resp.Status(), sku)
	}
}
