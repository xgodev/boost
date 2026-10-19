// Package boot is the smallest boost app: one annotated handler, one
// service, and the main builder generates, which runs boost.Start (config,
// then the zerolog logger) before the app. The config file sets the logger
// to JSON at DEBUG; GET /greetings/:name logs the greeting and answers it.
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl localhost:8080/greetings/world
package boot

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
