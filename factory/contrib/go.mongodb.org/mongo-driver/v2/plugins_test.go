package mongo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// unreachable makes the connection attempt fail fast after the client options
// plugins ran; the client plugins need a live server and are not covered here.
const unreachable = "mongodb://127.0.0.1:1/test"

func optionsPlugin(fn ClientOptionsPlugin) Plugin {
	return func(context.Context) (ClientOptionsPlugin, ClientPlugin) { return fn, nil }
}

func TestNewConnWithOptionsRunsClientOptionsPluginsInOrderOnTheClientOptions(t *testing.T) {
	var calls []string
	record := func(name string) Plugin {
		return optionsPlugin(func(_ context.Context, co *options.ClientOptions) error {
			calls = append(calls, name)
			co.SetServerSelectionTimeout(50 * time.Millisecond)
			return nil
		})
	}

	_, err := NewConnWithOptions(context.Background(), &Options{Uri: unreachable}, record("a"), record("b"))

	assert.Error(t, err, "no server listens on the unreachable address")
	assert.Equal(t, []string{"a", "b"}, calls)
}

func TestNewConnWithOptionsReturnsClientOptionsPluginError(t *testing.T) {
	boom := errors.New("boom")
	var after bool
	failing := optionsPlugin(func(context.Context, *options.ClientOptions) error { return boom })
	next := optionsPlugin(func(context.Context, *options.ClientOptions) error { after = true; return nil })

	conn, err := NewConnWithOptions(context.Background(), &Options{Uri: unreachable}, failing, next)

	assert.ErrorIs(t, err, boom)
	assert.ErrorContains(t, err, "failed to apply client options plugin")
	assert.Nil(t, conn)
	assert.False(t, after, "plugins after the failing one must not run")
}
