package resty

import (
	"context"
	"errors"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientWithOptionsRunsPluginsInOrderOnTheBuiltClient(t *testing.T) {
	var calls []string
	var got *resty.Client
	plugin := func(name string) Plugin {
		return func(_ context.Context, c *resty.Client) error {
			calls = append(calls, name)
			got = c
			return nil
		}
	}

	client := NewClientWithOptions(context.Background(), &Options{Host: "http://example"}, plugin("a"), plugin("b"))

	assert.Equal(t, []string{"a", "b"}, calls)
	require.NotNil(t, client)
	assert.Same(t, client, got)
}

func TestNewClientWithOptionsPanicsOnPluginError(t *testing.T) {
	boom := errors.New("boom")
	var after bool
	failing := func(context.Context, *resty.Client) error { return boom }
	next := func(context.Context, *resty.Client) error { after = true; return nil }

	assert.PanicsWithError(t, boom.Error(), func() {
		NewClientWithOptions(context.Background(), &Options{}, failing, next)
	})
	assert.False(t, after, "plugins after the failing one must not run")
}
