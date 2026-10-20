// Package rest holds what the rest server of app api needs besides its
// handlers: errorHandler turns boost errors into HTTP statuses, so a
// Conflict answers 409.
// @Server(type=rest, plugins=[recover, errorHandler])
package rest
