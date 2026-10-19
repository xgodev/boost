// Package service serves GET /reports/:id. The app provides a middleware
// wrapper, boost's cache middleware over a freecache manager then the native
// fallback middleware, and the report finder runs its calls through it: the
// first call for an id builds the report, the next ones answer from the
// cache, and a call that fails goes through the fallback and is not cached.
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl localhost:8080/reports/1
package service

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
