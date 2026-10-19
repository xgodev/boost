// Package middlewares runs a function behind the prometheus, logger and
// recovery middlewares of boost's function bootstrap, through the cloudevents
// http adapter (port 8080): an event without number panics, recovery turns
// the panic into an error, the logger logs it and prometheus counts it.
//
// configs/config.yaml makes the adapter the default command:
//
//	CONF=configs/config.yaml go run ./cmd/worker
//
// README.md posts events to it.
package middlewares

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
