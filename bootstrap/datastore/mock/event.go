package mock

import "github.com/xgodev/boost/bootstrap/repository"

// NewEvent returns a initialized mock that implements an event repository.
func NewEvent() repository.Event {
	return NewMock()
}
