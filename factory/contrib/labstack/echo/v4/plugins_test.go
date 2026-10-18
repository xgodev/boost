package echo

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServerWithOptionsRunsPluginsInOrderOnTheBuiltServer(t *testing.T) {
	var calls []string
	var got *Server
	plugin := func(name string) Plugin {
		return func(_ context.Context, s *Server) error {
			calls = append(calls, name)
			got = s
			return nil
		}
	}

	srv := NewServerWithOptions(context.Background(), &Options{}, plugin("a"), plugin("b"))

	assert.Equal(t, []string{"a", "b"}, calls)
	require.NotNil(t, srv)
	assert.Same(t, srv, got)
}

func TestNewServerWithOptionsPanicsOnPluginError(t *testing.T) {
	boom := errors.New("boom")
	var after bool
	failing := func(context.Context, *Server) error { return boom }
	next := func(context.Context, *Server) error { after = true; return nil }

	assert.PanicsWithError(t, boom.Error(), func() {
		NewServerWithOptions(context.Background(), &Options{}, failing, next)
	})
	assert.False(t, after, "plugins after the failing one must not run")
}
