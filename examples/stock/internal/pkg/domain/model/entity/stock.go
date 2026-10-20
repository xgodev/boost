// Package entity holds the entities of the domain: no transport, no driver
// tags.
package entity

import (
	"regexp"

	"github.com/xgodev/boost/model/errors"
)

// Stock is the units in stock of one item.
type Stock struct {
	SKU   string
	Units int64
}

var sku = regexp.MustCompile(`^[a-z0-9-]{1,64}$`)

// ValidSKU tells whether s is a sku: 1 to 64 of a-z, 0-9 and -. Every door
// reaches the stock through a service that checks it first.
func ValidSKU(s string) error {
	if !sku.MatchString(s) {
		return errors.NewNotValid(nil, s+" is not a sku: use 1 to 64 of a-z, 0-9 and -")
	}
	return nil
}
