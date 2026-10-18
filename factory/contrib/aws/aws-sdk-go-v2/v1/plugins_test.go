package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigWithOptionsRunsPluginsInOrderOnTheReturnedConfig(t *testing.T) {
	var calls []string
	plugin := func(name string) Plugin {
		return func(_ context.Context, cfg *aws.Config) error {
			calls = append(calls, name)
			cfg.Region = name
			return nil
		}
	}

	cfg, err := NewConfigWithOptions(context.Background(), &Options{DefaultRegion: "us-east-1"}, plugin("a"), plugin("b"))

	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, calls)
	assert.Equal(t, "b", cfg.Region, "changes made by plugins must reach the returned config")
}

func TestNewConfigWithOptionsReturnsPluginError(t *testing.T) {
	boom := errors.New("boom")
	var after bool
	failing := func(context.Context, *aws.Config) error { return boom }
	next := func(context.Context, *aws.Config) error { after = true; return nil }

	cfg, err := NewConfigWithOptions(context.Background(), &Options{DefaultRegion: "us-east-1"}, failing, next)

	assert.ErrorIs(t, err, boom)
	assert.Equal(t, aws.Config{}, cfg)
	assert.False(t, after, "plugins after the failing one must not run")
}
