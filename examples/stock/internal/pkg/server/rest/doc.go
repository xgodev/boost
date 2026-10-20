// Package rest holds what the rest server of app api is built with besides
// its handlers: recover from panics, a request id per call, a log line per
// call, boost errors as HTTP statuses (NotFound 404, Conflict 409), GET
// /health with the checkers the clients register and GET /metrics. Each
// plugin has its own enabled key under boost.factory.echo.plugins.
//
// @Server(type=rest, plugins=[recover, requestId, logger, errorHandler, health, prometheus])
package rest
