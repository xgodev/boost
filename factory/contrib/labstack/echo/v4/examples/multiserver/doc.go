// Package multiserver serves GET /test on echo next to a server of its own
// that is not echo. multiserver runs both, and GET /check answers 503 once
// the background server stops (10s after the start).
package multiserver

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
