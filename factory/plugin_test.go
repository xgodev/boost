package factory_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/factory"
)

type target struct {
	calls []string
}

type ctxKey struct{}

func record(name string) factory.Plugin[*target] {
	return func(_ context.Context, t *target) error {
		t.calls = append(t.calls, name)
		return nil
	}
}

func TestApplyPluginsRunsEveryPluginInOrder(t *testing.T) {
	tg := &target{}

	err := factory.ApplyPlugins(context.Background(), tg, record("a"), record("b"), record("c"))

	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, tg.calls)
}

func TestApplyPluginsStopsAtFirstError(t *testing.T) {
	tg := &target{}
	boom := errors.New("boom")
	failing := func(_ context.Context, t *target) error {
		t.calls = append(t.calls, "failing")
		return boom
	}

	err := factory.ApplyPlugins(context.Background(), tg, record("a"), failing, record("c"))

	assert.ErrorIs(t, err, boom)
	assert.Equal(t, []string{"a", "failing"}, tg.calls)
}

func TestApplyPluginsPassesContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey{}, "v")
	var got any
	p := func(ctx context.Context, _ *target) error {
		got = ctx.Value(ctxKey{})
		return nil
	}

	assert.NoError(t, factory.ApplyPlugins(ctx, &target{}, p))
	assert.Equal(t, "v", got)
}

func TestApplyPluginsWithoutPluginsIsNoop(t *testing.T) {
	assert.NoError(t, factory.ApplyPlugins[*target](context.Background(), &target{}))
}
