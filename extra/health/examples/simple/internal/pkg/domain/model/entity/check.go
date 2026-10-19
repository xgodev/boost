// Package entity holds the entities of the domain.
package entity

// Check is the result of a health checker.
type Check struct {
	Name     string
	Required bool
	// Error is why the check failed; empty when it passed.
	Error string
}
