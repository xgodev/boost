// Package service serves GET /greetings/:name. The app provides a middleware
// wrapper, boost's log middleware then a hystrix circuit breaker, and the
// greeter gets it by injection and runs its calls through it.
//
//	go run ./cmd/api
//	curl localhost:8080/greetings/ana
package service

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
