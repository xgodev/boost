// Package nats runs a function on the events published to nats subject
// messages, through the cloudevents nats adapter of boost's function bootstrap
// (boost.bootstrap.function.adapter.cloudevents.nats), which reads
// structured cloudevents.
//
// configs/config.yaml makes the adapter the default command and subscribes
// it to subject messages:
//
//	CONF=configs/config.yaml go run ./cmd/worker
//
// README.md starts nats and publishes an event to it.
package nats

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
