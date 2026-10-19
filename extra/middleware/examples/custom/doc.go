// Package custom serves GET /greetings/:name. The greeter runs its calls
// through a custom middleware, counter, which counts them: a middleware is a
// type with the Exec of middleware.AnyErrorMiddleware.
//
//	go run ./cmd/api
//	curl localhost:8080/greetings/ana
package custom

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
