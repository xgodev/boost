// Package envfile serves GET /application, the names of the application
// from wrapper/config, where a value comes from the first source that has
// it: an env var, then the file CONF names, then the default of config.Add.
// The file sets app.application.name; app.application.myName, a camelCase
// key, takes its env var with __ between the words: APP_APPLICATION_MY__NAME.
//
//	CONF=configs/config.yaml go run ./cmd/api
//	CONF=configs/config.yaml APP_APPLICATION_MY__NAME=mine go run ./cmd/api
//	curl localhost:8080/application
package envfile

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
