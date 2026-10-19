// Package consumer runs a function on the messages of a kafka topic, through
// the kafka adapter of boost's function bootstrap and the consumer of boost's
// factory (boost.factory.confluent).
//
// configs/config.yaml subscribes the topic messages and makes the kafka
// adapter the default command. Run it next to a broker
// (boost.factory.confluent.brokers, default localhost:9092):
//
//	CONF=configs/config.yaml go run ./cmd/worker
package consumer

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
