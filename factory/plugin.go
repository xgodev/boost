// Package factory holds the contract shared by every factory under it.
package factory

import "context"

// Plugin customizes a component T built by a factory (a client, a server,
// a config). Factories accept plugins as a variadic argument and run them
// right after building the component.
type Plugin[T any] func(context.Context, T) error

// ApplyPlugins runs plugins against target in order and stops at the first
// error, returning it.
func ApplyPlugins[T any](ctx context.Context, target T, plugins ...Plugin[T]) error {
	for _, plugin := range plugins {
		if err := plugin(ctx, target); err != nil {
			return err
		}
	}
	return nil
}
