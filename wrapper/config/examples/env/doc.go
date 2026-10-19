// Package env serves GET /application, the name of the application from
// wrapper/config: the key app.application.name has a default, an env var
// overrides it, and a config.Load after the env var changed shows in the
// next answer.
//
//	go run ./cmd/api
//	APP_APPLICATION_NAME=orders go run ./cmd/api
//	curl localhost:8080/application
package env

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
