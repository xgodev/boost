package fx

import (
	"go.uber.org/fx"
)

// NewApp start a new fx app with options, logging fx's events with
// NewEventLogger.
func NewApp(opts ...fx.Option) *fx.App {
	opts = append([]fx.Option{fx.WithLogger(NewEventLogger)}, opts...)
	return fx.New(opts...)
}
