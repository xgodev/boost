// Package recovery runs a function behind the recovery middleware of boost's
// function bootstrap, through the cloudevents http adapter (port 8080): an
// event without number panics, and recovery turns the panic into an error
// the adapter answers with 500, so the worker keeps running.
//
// configs/config.yaml makes the adapter the default command:
//
//	CONF=configs/config.yaml go run ./cmd/worker
//
// README.md posts events to it.
package recovery

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
