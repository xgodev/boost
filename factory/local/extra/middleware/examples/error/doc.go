// Package error serves GET /greetings/:name?lang=. The greeter runs its calls
// through the wrapper boost's middleware factory builds from config, the log
// middleware then a hystrix circuit breaker: a language it does not know
// fails, and the fallback of the circuit breaker greets in English.
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl "localhost:8080/greetings/ana?lang=pt"
//	curl "localhost:8080/greetings/ana?lang=xx"
package error

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
