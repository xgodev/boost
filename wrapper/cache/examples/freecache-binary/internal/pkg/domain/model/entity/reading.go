// Package entity holds the entities of the domain.
package entity

// Reading is what a sensor measured, where and when.
type Reading struct {
	Sensor      uint32
	Location    uint16
	Timestamp   uint32
	Temperature int16
}
