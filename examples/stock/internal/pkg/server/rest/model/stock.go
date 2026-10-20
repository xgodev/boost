// Package model holds the JSON bodies of the rest server.
package model

import "github.com/xgodev/boost/examples/stock/internal/pkg/domain/model/entity"

// Stock is the body of a stock: the request of POST /stocks and the
// response of every route.
type Stock struct {
	SKU   string `json:"sku"`
	Units int64  `json:"units"`
}

// FromEntity returns the body of the stock s.
func FromEntity(s *entity.Stock) Stock {
	return Stock{SKU: s.SKU, Units: s.Units}
}

// Entity returns the stock of the body.
func (s Stock) Entity() *entity.Stock {
	return &entity.Stock{SKU: s.SKU, Units: s.Units}
}
