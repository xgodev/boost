// Package simple serves GET /greetings/:name. The greeter runs its calls
// through the wrapper boost's middleware factory builds from config: the log
// middleware, then a hystrix circuit breaker, both named greeter in
// configs/config.yaml.
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl localhost:8080/greetings/ana
package simple

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
