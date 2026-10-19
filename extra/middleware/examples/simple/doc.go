// Package simple serves GET /greetings/:name. The greeter runs its calls
// through a middleware wrapper it builds itself: boost's log middleware, then
// a hystrix circuit breaker with a fallback.
//
//	go run ./cmd/api
//	curl localhost:8080/greetings/ana
package simple

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
