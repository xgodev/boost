// Package entity holds the entities of the domain.
package entity

import "time"

// File is a file the app keeps.
type File struct {
	Name         string
	LastModified time.Time
}
