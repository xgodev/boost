// Package http runs a function on the cloudevents posted to it, through the
// cloudevents http adapter of boost's function bootstrap
// (boost.bootstrap.function.adapter.cloudevents.http, port 8080).
//
// configs/config.yaml makes the adapter the default command:
//
//	CONF=configs/config.yaml go run ./cmd/worker
//
// README.md posts an event to it.
package http

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
