// Package entity holds the entities of the domain: no transport, no driver
// tags.
package entity

// Stock is the units in stock of one item.
type Stock struct {
	SKU   string
	Units int64
}
