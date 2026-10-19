// Package entity holds the entities of the domain.
package entity

// Database is the database the app talks to, as its last ping saw it.
type Database struct {
	Name   string
	Status string
}
