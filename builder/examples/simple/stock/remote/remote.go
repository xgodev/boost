// Package remote asks a warehouse over HTTP for the stock. Its client is the
// resty client builder generates for the name warehouse: one per service,
// with its own config (boost.factory.resty.clients.warehouse.*).
package remote

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/builder/examples/simple/stock"
	"github.com/xgodev/boost/model/errors"
)

type service struct {
	warehouse *resty.Client
}

type item struct {
	Available bool `json:"available"`
}

// New provides a stock service that asks the warehouse whether an item is
// available (GET /items/{item}); builder builds its client with the log
// and otel plugins, and the module of the client starts the OpenTelemetry
// providers.
// @Provide(index=0)
// @Inject(index=0, name=warehouse, plugins=[log, otel])
func New(warehouse *resty.Client) stock.Service {
	return &service{warehouse: warehouse}
}

func (s *service) Available(ctx context.Context, name string) (bool, error) {
	var got item
	resp, err := s.warehouse.R().
		SetContext(ctx).
		SetPathParam("item", name).
		SetResult(&got).
		Get("/items/{item}")
	if err != nil {
		return false, errors.NewInternal(err, "asking the warehouse for "+name)
	}
	if resp.StatusCode() != http.StatusOK {
		return false, errors.Internalf("the warehouse answered %s for %s", resp.Status(), name)
	}
	return got.Available, nil
}
