// Package entity holds the entities of the domain.
package entity

// Report is a report the app builds. Build is the number of the build that
// made it, so two answers with the same Build came from one build.
type Report struct {
	ID    string
	Build int64
}
