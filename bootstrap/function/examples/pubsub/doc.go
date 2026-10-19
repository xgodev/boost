// Package pubsub runs a function on the messages of a Google Cloud Pub/Sub
// subscription, through the pubsub adapter of boost's function bootstrap
// (boost.bootstrap.function.adapter.pubsub), and publishes its replies with
// the pubsub publisher driver.
//
// configs/config.yaml makes the adapter the default command, subscribes it
// to messages-worker and points the pubsub factory at an emulator:
//
//	CONF=configs/config.yaml go run ./cmd/worker
//
// README.md starts the emulator and publishes a message to it.
package pubsub

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
